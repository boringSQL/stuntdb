package facts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

const (
	Format = 1

	MetaFile     = "meta.json"
	StatsFile    = "stats.sql"
	SettingsFile = "settings.json"
	SchemaFile   = "schema.sql"
	WorkloadFile = "workload.json"
)

var maskedStates = []string{"no", "yes", "partly"}

type (
	Meta struct {
		Format   int       `json:"format"`
		Server   string    `json:"server"`
		Major    int       `json:"major"`
		Captured time.Time `json:"captured"`
		Hash     string    `json:"hash"`
		Masked   string    `json:"masked"`
		Producer string    `json:"producer"`
		PgDump   string    `json:"pg_dump"`
		Omitted  []string  `json:"omitted"`
	}

	Facts struct {
		Dir      string
		Meta     Meta
		Stats    []byte
		Settings []byte
		Schema   []byte
	}
)

func Open(dir string) (*Facts, error) {
	raw, err := os.ReadFile(filepath.Join(dir, MetaFile))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dir, err)
	}
	f := &Facts{Dir: dir}
	if err := json.Unmarshal(raw, &f.Meta); err != nil {
		return nil, fmt.Errorf("%s: %s: %w", dir, MetaFile, err)
	}
	for _, name := range []string{StatsFile, SettingsFile} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("%s: %s is required: %w", dir, name, err)
		}
		if name == StatsFile {
			f.Stats = data
		} else {
			f.Settings = data
		}
	}
	schema, err := os.ReadFile(filepath.Join(dir, SchemaFile))
	switch {
	case err == nil:
		f.Schema = schema
	case !errors.Is(err, fs.ErrNotExist):
		return nil, fmt.Errorf("%s: %w", dir, err)
	}
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return f, nil
}

func Write(dir string, f *Facts) error {
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s already exists", dir)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	f.Dir = dir
	f.Meta.Format = Format
	f.Meta.Hash = f.Hash()
	if f.Meta.Omitted == nil {
		f.Meta.Omitted = []string{}
	}
	if err := f.Validate(); err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(filepath.Dir(dir), ".facts-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := writeFiles(tmp, f); err != nil {
		return err
	}
	return os.Rename(tmp, dir)
}

func (f *Facts) Validate() error {
	if f.Meta.Format != Format {
		return fmt.Errorf("%s: format %d, want %d", f.Dir, f.Meta.Format, Format)
	}
	if !slices.Contains(maskedStates, f.Meta.Masked) {
		return fmt.Errorf("%s: masked %q, want no, yes or partly", f.Dir, f.Meta.Masked)
	}
	if f.Meta.Server == "" || f.Meta.Major == 0 {
		return fmt.Errorf("%s: %s names no server", f.Dir, MetaFile)
	}
	if f.Meta.Captured.IsZero() {
		return fmt.Errorf("%s: %s has no captured time", f.Dir, MetaFile)
	}
	if len(f.Stats) == 0 {
		return fmt.Errorf("%s: %s is empty", f.Dir, StatsFile)
	}
	if len(f.Settings) == 0 {
		return fmt.Errorf("%s: %s is empty", f.Dir, SettingsFile)
	}
	if f.Meta.Hash != f.Hash() {
		return fmt.Errorf("%s: hash does not match the files, capture again", f.Dir)
	}
	return nil
}

func (f *Facts) Hash() string {
	manifest := sha256.New()
	files := []struct {
		name string
		data []byte
	}{
		{StatsFile, f.Stats},
		{SettingsFile, f.Settings},
		{SchemaFile, f.Schema},
	}
	for _, file := range files {
		if file.data == nil {
			continue
		}
		fmt.Fprintf(manifest, "%x  %s\n", sha256.Sum256(file.data), file.name)
	}
	return "sha256:" + hex.EncodeToString(manifest.Sum(nil))
}

func writeFiles(dir string, f *Facts) error {
	meta, err := json.MarshalIndent(f.Meta, "", "  ")
	if err != nil {
		return err
	}
	files := []struct {
		name string
		data []byte
	}{
		{MetaFile, append(meta, '\n')},
		{StatsFile, f.Stats},
		{SettingsFile, f.Settings},
	}
	if f.Schema != nil {
		files = append(files, struct {
			name string
			data []byte
		}{SchemaFile, f.Schema})
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(dir, file.name), file.data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
