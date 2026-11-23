package stuntdb_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"testing"
	"time"

	"github.com/boringsql/stuntdb/internal/testpg"
)

const envHold = "STUNTDB_TEST_HOLD"

func TestCreateAndDrop(t *testing.T) {
	server := open(t)
	db, err := server.Create(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(db.Name(), "stuntdb_") {
		t.Errorf("name = %q", db.Name())
	}
	if !databaseExists(t, db.Name()) {
		t.Fatalf("%s was not created", db.Name())
	}
	if err := db.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if databaseExists(t, db.Name()) {
		t.Errorf("%s survived its drop", db.Name())
	}
	if err := db.Drop(context.Background()); err != nil {
		t.Errorf("second drop: %v", err)
	}
}

func TestHold(t *testing.T) {
	if os.Getenv(envHold) == "" {
		t.Skip("helper process for TestInterrupt")
	}
	server := open(t)
	db, err := server.Create(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(db.Name())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	if err := db.Drop(context.WithoutCancel(ctx)); err != nil {
		t.Fatal(err)
	}
}

func TestInterrupt(t *testing.T) {
	if os.Getenv(envHold) != "" {
		t.Skip("helper process")
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestHold$")
	cmd.Env = append(os.Environ(), envHold+"=1")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	var finished bool
	t.Cleanup(func() {
		if finished {
			return
		}
		cmd.Process.Kill()
		<-done
	})

	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatalf("helper printed no name: %v", scanner.Err())
	}
	name := scanner.Text()
	if !databaseExists(t, name) {
		t.Fatalf("%s is not there", name)
	}

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		finished = true
		if err != nil {
			t.Fatalf("helper: %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("helper ignored the interrupt")
	}

	if databaseExists(t, name) {
		t.Errorf("%s survived the interrupt", name)
	}
}

func databaseExists(t *testing.T, name string) bool {
	t.Helper()
	db := testpg.Open(t, testpg.Admin(t))
	var found bool
	if err := db.QueryRow("select exists (select 1 from pg_database where datname = $1)", name).Scan(&found); err != nil {
		t.Fatal(err)
	}
	return found
}
