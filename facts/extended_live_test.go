package facts_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/facts"
	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestExtendedStatistics(t *testing.T) {
	ctx := context.Background()
	source := testpg.Create(t)
	db := testpg.Open(t, source)
	for _, stmt := range []string{
		"create table people(city text, country text)",
		"insert into people select 'c' || (g % 7), 'k' || (g % 3) from generate_series(1, 500) g",
		"create statistics people_deps (dependencies) on city, country from people",
		"analyze",
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}

	dir := filepath.Join(t.TempDir(), "extended.facts")
	captured, _, err := facts.Capture(ctx, source, dir, "stuntdb test")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(captured.Schema), "CREATE STATISTICS") {
		t.Error("schema.sql lost the extended statistics definition")
	}
	hasData := strings.Contains(string(captured.Stats), "pg_restore_extended_stats")
	if testpg.Major(t) < 19 {
		if hasData {
			t.Error("a statistics file before 19 carried extended statistics data")
		}
		if len(captured.Meta.Omitted) == 0 {
			t.Error("meta.json omitted nothing, but no extended statistics data was captured")
		}
		return
	}
	if !hasData {
		t.Error("19 carried no extended statistics data")
	}
	if len(captured.Meta.Omitted) != 0 {
		t.Errorf("omitted = %v", captured.Meta.Omitted)
	}
}
