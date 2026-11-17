package facts_test

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/boringsql/stuntdb/facts"
	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestNewMeta(t *testing.T) {
	db := testpg.Open(t, testpg.Create(t))
	meta, err := facts.NewMeta(db, "stuntdb test", "18.6")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Format != facts.Format {
		t.Errorf("format = %d", meta.Format)
	}
	if meta.Major != testpg.Major(t) {
		t.Errorf("major = %d, server says %d", meta.Major, testpg.Major(t))
	}
	if !strings.HasPrefix(meta.Server, strconv.Itoa(meta.Major)+".") {
		t.Errorf("server = %q", meta.Server)
	}
	if age := time.Since(meta.Captured); age > time.Minute || age < 0 {
		t.Errorf("captured = %s", meta.Captured)
	}
	if meta.Masked != "no" {
		t.Errorf("masked = %q", meta.Masked)
	}
	if meta.Producer != "stuntdb test" || meta.PgDump != "18.6" {
		t.Errorf("producer = %q, pg_dump = %q", meta.Producer, meta.PgDump)
	}
	if meta.Major < 19 && len(meta.Omitted) == 0 {
		t.Errorf("omitted nothing on %d", meta.Major)
	}
	if meta.Major >= 19 && len(meta.Omitted) != 0 {
		t.Errorf("omitted something on %d: %v", meta.Major, meta.Omitted)
	}
}

func TestWriteMeta(t *testing.T) {
	meta, err := facts.NewMeta(testpg.Open(t, testpg.Create(t)), "stuntdb test", "18.6")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "capture.facts")
	written := &facts.Facts{
		Meta:     meta,
		Stats:    []byte("select 1;\n"),
		Settings: []byte("{\"work_mem\":\"4MB\"}\n"),
	}
	if err := facts.Write(out, written); err != nil {
		t.Fatal(err)
	}
	back, err := facts.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	if back.Meta.Hash != back.Hash() {
		t.Errorf("meta says %s, files hash to %s", back.Meta.Hash, back.Hash())
	}
	if back.Meta.Server != meta.Server || !back.Meta.Captured.Equal(meta.Captured) {
		t.Errorf("server = %q, captured = %s", back.Meta.Server, back.Meta.Captured)
	}
}

func TestCaptureCounts(t *testing.T) {
	db := testpg.Open(t, testpg.Create(t))
	for _, stmt := range []string{
		"create table tickets(id bigint primary key, subject text, opened_at timestamptz)",
		"create index tickets_opened_idx on tickets(opened_at)",
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	counts, err := facts.CaptureCounts(db)
	if err != nil {
		t.Fatal(err)
	}
	want := facts.Counts{Tables: 1, Columns: 3, Indexes: 2, Relations: 3}
	if counts != want {
		t.Errorf("counts = %+v, want %+v", counts, want)
	}
}
