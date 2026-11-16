package facts_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/facts"
)

func TestLoadSettings(t *testing.T) {
	settings, err := facts.LoadSettings([]byte(`{"random_page_cost":"1.1","work_mem":"4MB"}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(settings) != 2 || settings["random_page_cost"] != "1.1" || settings["work_mem"] != "4MB" {
		t.Errorf("settings = %v", settings)
	}
}

func TestLoadSettingsUnknownName(t *testing.T) {
	_, err := facts.LoadSettings([]byte(`{"track_io_timing":"on"}`))
	if err == nil || !strings.Contains(err.Error(), `"track_io_timing"`) {
		t.Errorf("error = %v, want it to name track_io_timing", err)
	}
}

func TestLoadSettingsNotString(t *testing.T) {
	_, err := facts.LoadSettings([]byte(`{"work_mem":4}`))
	if err == nil || !strings.Contains(err.Error(), "must be a string") {
		t.Errorf("error = %v", err)
	}
}

func TestMarshalSettings(t *testing.T) {
	got, err := facts.MarshalSettings(map[string]string{"work_mem": "4MB", "random_page_cost": "1.1"})
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"random_page_cost\": \"1.1\",\n  \"work_mem\": \"4MB\"\n}\n"
	if string(got) != want {
		t.Errorf("settings.json = %q, want %q", got, want)
	}
}

func TestOpenRejectsForeignKey(t *testing.T) {
	dir := copyFixture(t)
	path := filepath.Join(dir, facts.SettingsFile)
	if err := os.WriteFile(path, []byte(`{"track_io_timing":"on"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	assertOpenError(t, dir, "track_io_timing")
}
