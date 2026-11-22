package pgdump_test

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/pgdump"
	"github.com/boringsql/stuntdb/internal/testpg"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestDump(t *testing.T) {
	ctx := context.Background()
	uri := testpg.Create(t)
	db := testpg.Open(t, uri)
	for _, stmt := range []string{
		"create table customers(id bigint primary key, email text)",
		"create table orders(id bigint primary key, customer_id bigint references customers, total numeric(12,2))",
		"create table events(id bigint primary key, kind text)",
		"insert into customers select g, 'user' || g from generate_series(1, 100) g",
		"insert into orders select g, 1 + g % 100, g from generate_series(1, 1000) g",
		"insert into events select g, 'kind' || (g % 5) from generate_series(1, 500) g",
		"analyze",
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}

	var source string
	if err := db.QueryRow("select current_setting('server_version')").Scan(&source); err != nil {
		t.Fatal(err)
	}

	tool, err := pgdump.Find()
	if err != nil {
		t.Fatal(err)
	}
	if err := tool.Check(source); err != nil {
		t.Fatal(err)
	}

	stats, err := tool.Statistics(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	if calls := statisticsCalls(t, stats); calls < 3 {
		t.Errorf("%d statistics calls for three tables", calls)
	}

	schema, err := tool.Schema(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	other := simpleDB(t, testpg.Create(t))
	if _, err := other.Exec(stripMetaCommands(string(schema))); err != nil {
		t.Fatalf("loading schema.sql: %v", err)
	}
	var tables int
	if err := other.QueryRow("select count(*) from pg_tables where schemaname = 'public'").Scan(&tables); err != nil {
		t.Fatal(err)
	}
	if tables != 3 {
		t.Errorf("%d tables after loading schema.sql, want 3", tables)
	}
}

func TestOlderClientRefusesNewerServer(t *testing.T) {
	tool, err := pgdump.Find()
	if err != nil {
		t.Fatal(err)
	}
	if tool.Version.Major >= 19 {
		t.Skipf("client %s can dump a 19 server", tool.Version)
	}
	err = tool.Check("19.1")
	if err == nil || !strings.Contains(err.Error(), "cannot dump a 19.1 server") {
		t.Errorf("check = %v", err)
	}
}

func statisticsCalls(t *testing.T, stats []byte) int {
	t.Helper()
	text := string(stats)
	if !strings.Contains(text, "\n\\restrict ") {
		t.Errorf("stats.sql has no \\restrict line")
	}
	if !strings.Contains(text, "\n\\unrestrict ") {
		t.Errorf("stats.sql has no \\unrestrict line")
	}
	calls, inside := 0, false
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "SELECT * FROM pg_catalog.pg_restore_"):
			if inside {
				t.Errorf("a second call starts before the first ends")
			}
			inside = true
			calls++
		case line == ");":
			if !inside {
				t.Errorf("); outside a call")
			}
			inside = false
		}
	}
	if inside {
		t.Errorf("the last call is unfinished")
	}
	return calls
}

func stripMetaCommands(sql string) string {
	var b strings.Builder
	for _, line := range strings.Split(sql, "\n") {
		if strings.HasPrefix(line, "\\") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func simpleDB(t *testing.T, uri string) *sql.DB {
	t.Helper()
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	query := u.Query()
	query.Set("default_query_exec_mode", "simple_protocol")
	u.RawQuery = query.Encode()
	db, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
