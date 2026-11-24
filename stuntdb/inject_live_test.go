package stuntdb_test

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestInjectLabStatistics(t *testing.T) {
	ctx := context.Background()
	db := created(t, open(t))
	if err := db.Load(ctx, labSchema(t)); err != nil {
		t.Fatal(err)
	}
	skipped, err := db.Inject(ctx, labStats(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v", skipped)
	}

	stunt := connect(t, db.Name())
	var reltuples float64
	if err := stunt.QueryRow("select reltuples from pg_class where relname = 'orders'").Scan(&reltuples); err != nil {
		t.Fatal(err)
	}
	if reltuples != 2000000 {
		t.Errorf("orders has %v rows", reltuples)
	}

	var options string
	if err := stunt.QueryRow("select coalesce(array_to_string(reloptions, ','), '') from pg_class where relname = 'orders'").Scan(&options); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(options, "autovacuum_enabled=false") {
		t.Errorf("orders options = %q", options)
	}

	var toast string
	if err := stunt.QueryRow(`select coalesce(array_to_string(t.reloptions, ','), '')
		from pg_class p join pg_class t on t.oid = p.reltoastrelid
		where p.relname = 'orders'`).Scan(&toast); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(toast, "autovacuum_enabled=false") {
		t.Errorf("orders toast options = %q", toast)
	}

	var plan string
	if err := stunt.QueryRow("explain select * from orders").Scan(&plan); err != nil {
		t.Fatal(err)
	}
	rows := regexp.MustCompile(`rows=(\d+)`).FindStringSubmatch(plan)
	if rows == nil || rows[1] != "2000000" {
		t.Errorf("plan = %q", plan)
	}
}

func TestInjectListsRelationsWithoutStatistics(t *testing.T) {
	ctx := context.Background()
	schema := append(labSchema(t), []byte("create table public.audit(id bigint);\n")...)
	db := created(t, open(t))
	if err := db.Load(ctx, schema); err != nil {
		t.Fatal(err)
	}
	skipped, err := db.Inject(ctx, labStats(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 1 {
		t.Fatalf("skipped = %v", skipped)
	}
	if skipped[0].Subject != "public.audit" || !strings.Contains(skipped[0].Reason, "no statistics") {
		t.Errorf("skipped = %v", skipped)
	}
}

func TestInjectRecordsUnapplied(t *testing.T) {
	ctx := context.Background()
	db := created(t, open(t))
	if err := db.Load(ctx, []byte("create table people(city text);\n")); err != nil {
		t.Fatal(err)
	}
	stats := []byte(`SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'people',
	'relpages', '1'::integer,
	'reltuples', '10'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'people',
	'attname', 'city',
	'inherited', 'f'::boolean,
	'most_common_vals', '{{a},{b}}'::text[][]
);
`)
	skipped, err := db.Inject(ctx, stats)
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 1 || skipped[0].Subject != "public.people.city" || skipped[0].Reason != "not applied" {
		t.Errorf("skipped = %v", skipped)
	}
}

func TestInjectSkipsMissingRelations(t *testing.T) {
	ctx := context.Background()
	db := created(t, open(t))
	if err := db.Load(ctx, nil); err != nil {
		t.Fatal(err)
	}
	skipped, err := db.Inject(ctx, labStats(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 16 {
		t.Fatalf("%d skipped: %v", len(skipped), skipped)
	}
	for _, entry := range skipped {
		if entry.Reason != "no such relation" {
			t.Errorf("skipped = %v", entry)
		}
	}
}

func labSchema(t *testing.T) []byte {
	t.Helper()
	schema, err := os.ReadFile("../testdata/lab/facts/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	return schema
}

func labStats(t *testing.T) []byte {
	t.Helper()
	stats, err := os.ReadFile("../testdata/lab/facts/stats.sql")
	if err != nil {
		t.Fatal(err)
	}
	return stats
}
