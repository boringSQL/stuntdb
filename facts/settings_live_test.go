package facts_test

import (
	"testing"

	"github.com/boringsql/stuntdb/facts"
	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestCaptureSettings(t *testing.T) {
	db := testpg.Open(t, testpg.Create(t))
	settings, err := facts.CaptureSettings(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(settings) < 40 {
		t.Errorf("captured %d settings", len(settings))
	}

	var shown string
	if err := db.QueryRow("select current_setting('random_page_cost')").Scan(&shown); err != nil {
		t.Fatal(err)
	}
	if settings["random_page_cost"] != shown {
		t.Errorf("random_page_cost = %q, show says %q", settings["random_page_cost"], shown)
	}

	data, err := facts.MarshalSettings(settings)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := facts.LoadSettings(data); err != nil {
		t.Errorf("the capture writes a name the loader refuses: %v", err)
	}

	again, err := facts.CaptureSettings(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != len(settings) {
		t.Errorf("second capture has %d settings, first had %d", len(again), len(settings))
	}
}
