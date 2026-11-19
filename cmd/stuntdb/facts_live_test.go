package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/facts"
	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestFactsTranscript(t *testing.T) {
	shift := shiftSource(t)
	dir := filepath.Join(t.TempDir(), "capture.facts")

	var stdout, stderr bytes.Buffer
	if code := runFacts([]string{"facts", shift, "-o", dir}, &stdout, &stderr); code != exitOK {
		t.Fatalf("exit code = %d: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("transcript has %d lines:\n%s", len(lines), stdout.String())
	}
	for i, prefix := range []string{"statistics  ", "settings    ", "schema      ", "wrote "} {
		if !strings.HasPrefix(lines[i], prefix) {
			t.Errorf("line %d = %q, want prefix %q", i+1, lines[i], prefix)
		}
	}
	if !strings.Contains(lines[0], "3 relations, 3 columns") {
		t.Errorf("statistics line = %q", lines[0])
	}
	if !strings.Contains(lines[1], "planner settings") {
		t.Errorf("settings line = %q", lines[1])
	}
	if !strings.Contains(lines[2], "1 tables, 2 indexes") {
		t.Errorf("schema line = %q", lines[2])
	}
	if !strings.Contains(lines[3], dir) || !strings.Contains(lines[3], "not masked") {
		t.Errorf("wrote line = %q", lines[3])
	}

	captured, err := facts.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if captured.Meta.Hash != captured.Hash() {
		t.Errorf("meta says %s, files hash to %s", captured.Meta.Hash, captured.Hash())
	}
	if captured.Meta.Masked != "no" {
		t.Errorf("masked = %q", captured.Meta.Masked)
	}
	if len(captured.Stats) == 0 || len(captured.Schema) == 0 {
		t.Errorf("stats %d bytes, schema %d bytes", len(captured.Stats), len(captured.Schema))
	}
}

func shiftSource(t *testing.T) string {
	t.Helper()
	uri := testpg.Create(t)
	db := testpg.Open(t, uri)
	for _, stmt := range []string{
		"create table shifts(id bigint primary key, role text, starts_at timestamptz)",
		"create index shifts_starts_idx on shifts(starts_at)",
		"insert into shifts select g, 'role' || (g % 4), now() from generate_series(1, 100) g",
		"analyze",
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return uri
}
