package facts_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/boringsql/stuntdb/facts"
)

const (
	fixtureDir  = "../testdata/facts"
	fixtureHash = "sha256:768ec7d59bab474d2273c3720d14a8ddbd7b2f5e453808342e53438ca973a6a1"
)

func TestOpen(t *testing.T) {
	f, err := facts.Open(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	if f.Meta.Format != facts.Format {
		t.Errorf("format = %d", f.Meta.Format)
	}
	if f.Meta.Server != "18.1" || f.Meta.Major != 18 {
		t.Errorf("server = %s, major = %d", f.Meta.Server, f.Meta.Major)
	}
	if !f.Meta.Captured.Equal(time.Date(2026, 11, 22, 9, 0, 0, 0, time.UTC)) {
		t.Errorf("captured = %s", f.Meta.Captured)
	}
	if f.Meta.Masked != "no" {
		t.Errorf("masked = %q", f.Meta.Masked)
	}
	if f.Meta.Hash != fixtureHash || f.Hash() != fixtureHash {
		t.Errorf("hash = %s, want %s", f.Hash(), fixtureHash)
	}
	if len(f.Stats) == 0 || len(f.Settings) == 0 || len(f.Schema) == 0 {
		t.Errorf("stats %d bytes, settings %d, schema %d", len(f.Stats), len(f.Settings), len(f.Schema))
	}
}

func TestRoundTrip(t *testing.T) {
	orig, err := facts.Open(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "copy.facts")
	written := &facts.Facts{Meta: orig.Meta, Stats: orig.Stats, Settings: orig.Settings, Schema: orig.Schema}
	if err := facts.Write(out, written); err != nil {
		t.Fatal(err)
	}
	back, err := facts.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(back.Stats, orig.Stats) || !bytes.Equal(back.Settings, orig.Settings) || !bytes.Equal(back.Schema, orig.Schema) {
		t.Error("the files changed on the way through")
	}
	if back.Hash() != orig.Hash() {
		t.Errorf("hash = %s, want %s", back.Hash(), orig.Hash())
	}
	for _, name := range []string{facts.MetaFile, facts.StatsFile, facts.SettingsFile, facts.SchemaFile} {
		want, err := os.ReadFile(filepath.Join(fixtureDir, name))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s is not the file it was written from:\n%s", name, got)
		}
	}
	if _, err := os.Stat(filepath.Join(out, facts.WorkloadFile)); !os.IsNotExist(err) {
		t.Errorf("%s was copied, but it is not ours", facts.WorkloadFile)
	}
}

func TestNoSchema(t *testing.T) {
	orig, err := facts.Open(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "noschema.facts")
	written := &facts.Facts{Meta: orig.Meta, Stats: orig.Stats, Settings: orig.Settings}
	if err := facts.Write(out, written); err != nil {
		t.Fatal(err)
	}
	back, err := facts.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	if back.Schema != nil {
		t.Errorf("schema = %d bytes, want none", len(back.Schema))
	}
	if back.Hash() == orig.Hash() {
		t.Error("the hash did not change when schema.sql was left out")
	}
	if back.Meta.Hash != back.Hash() {
		t.Errorf("meta says %s, files hash to %s", back.Meta.Hash, back.Hash())
	}
}

func TestMissingStats(t *testing.T) {
	dir := copyFixture(t)
	if err := os.Remove(filepath.Join(dir, facts.StatsFile)); err != nil {
		t.Fatal(err)
	}
	assertOpenError(t, dir, "stats.sql")
}

func TestUnknownFormat(t *testing.T) {
	dir := copyFixture(t)
	path := filepath.Join(dir, facts.MetaFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var meta facts.Meta
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatal(err)
	}
	meta.Format = 2
	raw, err = json.MarshalIndent(meta, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	assertOpenError(t, dir, "format 2")
}

func TestStaleHash(t *testing.T) {
	dir := copyFixture(t)
	path := filepath.Join(dir, facts.StatsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	assertOpenError(t, dir, "hash does not match")
}

func TestExistingDir(t *testing.T) {
	dir := copyFixture(t)
	f, err := facts.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := facts.Write(dir, f); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("writing into an existing directory: %v", err)
	}
}

func copyFixture(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "facts")
	if err := os.CopyFS(dir, os.DirFS(fixtureDir)); err != nil {
		t.Fatal(err)
	}
	return dir
}

func assertOpenError(t *testing.T, dir, want string) {
	t.Helper()
	_, err := facts.Open(dir)
	if err == nil {
		t.Fatalf("opened %s, want %q", dir, want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want it to name %q", err, want)
	}
}
