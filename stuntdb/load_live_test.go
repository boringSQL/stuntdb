package stuntdb_test

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/testpg"
	"github.com/boringsql/stuntdb/stuntdb"
)

func TestLoad(t *testing.T) {
	db := created(t, open(t))
	schema := []byte("create table rides(id bigint primary key, zone text);\ncreate index rides_zone_idx on rides(zone);\n")
	if err := db.Load(context.Background(), schema); err != nil {
		t.Fatal(err)
	}

	stunt := connect(t, db.Name())
	if tables := relations(t, stunt, "r"); tables != 1 {
		t.Errorf("tables = %d", tables)
	}
	if indexes := relations(t, stunt, "i"); indexes != 2 {
		t.Errorf("indexes = %d, want 2", indexes)
	}

	var extensions int
	if err := stunt.QueryRow("select count(*) from pg_extension where extname = 'pg_regresql'").Scan(&extensions); err != nil {
		t.Fatal(err)
	}
	if extensions != 1 {
		t.Errorf("pg_regresql extensions = %d", extensions)
	}

	var active *string
	if err := stunt.QueryRow("select current_setting('pg_regresql.active', true)").Scan(&active); err != nil {
		t.Fatal(err)
	}
	if active == nil || *active != "on" {
		t.Errorf("pg_regresql.active = %v", active)
	}
}

func TestLoadLabSchema(t *testing.T) {
	db := created(t, open(t))
	if err := db.Load(context.Background(), labSchema(t)); err != nil {
		t.Fatal(err)
	}

	stunt := connect(t, db.Name())
	if tables := relations(t, stunt, "r"); tables != 5 {
		t.Errorf("tables = %d, want 5", tables)
	}
	if indexes := relations(t, stunt, "i"); indexes != 11 {
		t.Errorf("indexes = %d, want 11", indexes)
	}
}

func TestLoadRejectsBadSQL(t *testing.T) {
	db := created(t, open(t))
	err := db.Load(context.Background(), []byte("create table before_error(id int);\nnonsense;\n"))
	if err == nil {
		t.Fatal("loaded invalid sql")
	}
	if !strings.Contains(err.Error(), "syntax error") {
		t.Errorf("error = %q", err)
	}
}

func created(t *testing.T, server *stuntdb.Server) *stuntdb.Database {
	t.Helper()
	db, err := server.Create(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Drop(context.Background()); err != nil {
			t.Errorf("drop %s: %v", db.Name(), err)
		}
	})
	return db
}

func connect(t *testing.T, name string) *sql.DB {
	t.Helper()
	u, err := url.Parse(testpg.Admin(t))
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/" + name
	return testpg.Open(t, u.String())
}

func relations(t *testing.T, db *sql.DB, kind string) int {
	t.Helper()
	var n int
	err := db.QueryRow(`select count(*)
	from pg_class c join pg_namespace n on n.oid = c.relnamespace
	where n.nspname = 'public' and c.relkind = $1`, kind).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	return n
}
