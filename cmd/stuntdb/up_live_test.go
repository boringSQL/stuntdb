package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/boringsql/stuntdb/facts"
	"github.com/boringsql/stuntdb/internal/testpg"
)

const (
	labFacts  = "../../testdata/lab/facts"
	envHelper = "STUNTDB_TEST_UP"
)

func TestUpTranscript(t *testing.T) {
	server := testpg.Admin(t)
	var stdout, stderr bytes.Buffer
	code := runUp([]string{"up", labFacts, "--server", server, "--",
		"sh", "-c", `psql "$DATABASE_URL" -Atc 'select current_database()'`}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("exit code = %d: %s", code, stderr.String())
	}

	captured, err := facts.Open(labFacts)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) < 4 {
		t.Fatalf("transcript has %d lines:\n%s", len(lines), stdout.String())
	}
	name := databaseName(t, lines[0])
	want := fmt.Sprintf("plans like production as of %s. Statistics are not masked.",
		captured.Meta.Captured.Format("2006-01-02"))
	if lines[1] != want {
		t.Errorf("line 2 = %q, want %q", lines[1], want)
	}
	if want := "Empty tables: plan here, do not execute. Read-only. Ctrl-C drops it."; lines[2] != want {
		t.Errorf("line 3 = %q, want %q", lines[2], want)
	}
	for i, line := range lines[3 : len(lines)-1] {
		if !strings.HasPrefix(line, "skipped  ") {
			t.Errorf("line %d = %q, want a skipped line", i+4, line)
		}
	}
	if got := lines[len(lines)-1]; got != name {
		t.Errorf("command saw %q, transcript named %q", got, name)
	}
	assertGone(t, name)
}

func TestUpMasked(t *testing.T) {
	server := testpg.Admin(t)
	captured, err := facts.Open(labFacts)
	if err != nil {
		t.Fatal(err)
	}
	for masked, line := range map[string]string{
		"yes":    "Statistics are masked.",
		"partly": "Statistics are partly masked.",
	} {
		t.Run(masked, func(t *testing.T) {
			copied := *captured
			copied.Meta.Masked = masked
			dir := filepath.Join(t.TempDir(), "masked.facts")
			if err := facts.Write(dir, &copied); err != nil {
				t.Fatal(err)
			}

			var stdout, stderr bytes.Buffer
			if code := runUp([]string{"up", dir, "--server", server, "--", "true"}, &stdout, &stderr); code != exitOK {
				t.Fatalf("exit code = %d: %s", code, stderr.String())
			}
			want := fmt.Sprintf("plans like production as of %s. %s",
				captured.Meta.Captured.Format("2006-01-02"), line)
			if lines := strings.Split(stdout.String(), "\n"); lines[1] != want {
				t.Errorf("line 2 = %q, want %q", lines[1], want)
			}
		})
	}
}

func TestUpPlansLikeProduction(t *testing.T) {
	server := testpg.Admin(t)
	var stdout, stderr bytes.Buffer
	code := runUp([]string{"up", labFacts, "--server", server, "--",
		"sh", "-c", `psql "$DATABASE_URL" -Atc 'explain select * from orders'`}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("exit code = %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "rows=2000000") {
		t.Errorf("plan is not production's:\n%s", stdout.String())
	}
}

func TestUpSchemaOverride(t *testing.T) {
	schema := filepath.Join(t.TempDir(), "schema.sql")
	if err := os.WriteFile(schema, []byte("create table overridden(id int primary key);\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := testpg.Admin(t)
	var stdout, stderr bytes.Buffer
	code := runUp([]string{"up", labFacts, "--server", server, "--schema", schema, "--",
		"sh", "-c", `psql "$DATABASE_URL" -Atc "select to_regclass('overridden') is not null, to_regclass('orders') is null"`}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("exit code = %d: %s", code, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); !strings.HasSuffix(got, "t|t") {
		t.Errorf("schema check = %q:\n%s", got, stdout.String())
	}
}

func TestUpCommandExitCode(t *testing.T) {
	server := testpg.Admin(t)
	var stdout, stderr bytes.Buffer
	code := runUp([]string{"up", labFacts, "--server", server, "--",
		"sh", "-c", `psql "$DATABASE_URL" -Atc 'select current_database()'; exit 7`}, &stdout, &stderr)
	if code != 7 {
		t.Fatalf("exit code = %d, want 7: %s", code, stderr.String())
	}
	name := databaseName(t, strings.Split(stdout.String(), "\n")[0])
	assertGone(t, name)
}

func TestUpHolds(t *testing.T) {
	server := testpg.Admin(t)
	cmd := exec.Command(os.Args[0], "-test.run=^TestUpHelper$")
	cmd.Env = append(os.Environ(), envHelper+"=1", "STUNTDB_TEST_SERVER="+server)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	var name string
	for i := 0; i < 3; i++ {
		if !scanner.Scan() {
			cmd.Process.Kill()
			t.Fatalf("stopped after %d lines: %v", i, scanner.Err())
		}
		if i == 0 {
			name = databaseName(t, scanner.Text())
		}
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		t.Fatalf("up exited without an interrupt: %v", err)
	case <-time.After(250 * time.Millisecond):
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatalf("interrupted up: %v", err)
	}
	assertGone(t, name)
}

func TestUpHelper(t *testing.T) {
	if os.Getenv(envHelper) == "" {
		t.Skip("helper process for TestUpHolds")
	}
	args := []string{"up", labFacts, "--server", os.Getenv("STUNTDB_TEST_SERVER")}
	if code := runUp(args, os.Stdout, os.Stderr); code != exitOK {
		t.Fatalf("exit code = %d", code)
	}
}

func databaseName(t *testing.T, line string) string {
	t.Helper()
	uri, err := url.Parse(line)
	if err != nil {
		t.Fatalf("%q: %v", line, err)
	}
	name := path.Base(uri.Path)
	if !strings.HasPrefix(name, "stuntdb_") {
		t.Fatalf("uri %q names no stunt database", line)
	}
	if uri.User.Username() != name {
		t.Errorf("role %q for database %q", uri.User.Username(), name)
	}
	if password, _ := uri.User.Password(); len(password) != 24 {
		t.Errorf("password %q is not 24 characters", password)
	}
	return name
}

func assertGone(t *testing.T, name string) {
	t.Helper()
	var database, role bool
	err := testpg.Open(t, testpg.Admin(t)).QueryRow(`select
		exists (select 1 from pg_database where datname = $1),
		exists (select 1 from pg_roles where rolname = $1)`, name).Scan(&database, &role)
	if err != nil {
		t.Fatal(err)
	}
	if database || role {
		t.Errorf("%s left behind: database %t, role %t", name, database, role)
	}
}
