package stuntdb_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestSettings(t *testing.T) {
	ctx := context.Background()
	db := created(t, open(t))
	settings := map[string]string{"random_page_cost": "1.5", "work_mem": "8192"}
	skipped, err := db.Settings(ctx, settings)
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v", skipped)
	}

	stunt := connect(t, db.Name())
	for name, want := range settings {
		var got string
		if err := stunt.QueryRow("select setting from pg_settings where name = $1", name).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}

func TestSettingsSkips(t *testing.T) {
	ctx := context.Background()
	db := created(t, open(t))
	skipped, err := db.Settings(ctx, map[string]string{
		"random_page_cost":    "1.5",
		"work_mem":            "999TB",
		"stuntdb_nonexistent": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 2 {
		t.Fatalf("skipped = %v", skipped)
	}
	if skipped[0].Subject != "stuntdb_nonexistent=1" || skipped[0].Reason != "unknown to this server" {
		t.Errorf("skipped[0] = %v", skipped[0])
	}
	if skipped[1].Subject != "work_mem=999TB" || !strings.Contains(skipped[1].Reason, "invalid value") {
		t.Errorf("skipped[1] = %v", skipped[1])
	}
}

func TestSettingsLabFacts(t *testing.T) {
	ctx := context.Background()
	settings := labSettings(t)
	db := created(t, open(t))
	skipped, err := db.Settings(ctx, settings)
	if err != nil {
		t.Fatal(err)
	}

	want := 0
	switch testpg.Major(t) {
	case 16:
		want = 3
	case 17:
		want = 2
	}
	if len(skipped) != want {
		t.Errorf("%d skipped, want %d: %v", len(skipped), want, skipped)
	}
	for _, entry := range skipped {
		if entry.Reason != "unknown to this server" {
			t.Errorf("skipped = %v", entry)
		}
	}

	stunt := connect(t, db.Name())
	for _, name := range []string{"random_page_cost", "work_mem", "plan_cache_mode"} {
		var got string
		if err := stunt.QueryRow("select setting from pg_settings where name = $1", name).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != settings[name] {
			t.Errorf("%s = %q, want %q", name, got, settings[name])
		}
	}
}

func labSettings(t *testing.T) map[string]string {
	t.Helper()
	data, err := os.ReadFile("../testdata/lab/facts/settings.json")
	if err != nil {
		t.Fatal(err)
	}
	var settings map[string]string
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatal(err)
	}
	return settings
}
