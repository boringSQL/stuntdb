package stuntdb_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestReadOnly(t *testing.T) {
	ctx := context.Background()
	db := created(t, open(t))
	if err := db.Load(ctx, labSchema(t)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Inject(ctx, labStats(t)); err != nil {
		t.Fatal(err)
	}

	owner := connect(t, db.Name())
	for _, stmt := range []string{
		"alter table orders enable row level security",
		"create policy orders_tenant on orders using (tenant_id = 1)",
	} {
		if _, err := owner.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}

	role, err := db.ReadOnly(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if role.Name() != db.Name() {
		t.Errorf("role %q for database %q", role.Name(), db.Name())
	}

	readOnly := testpg.Open(t, role.URI())
	var rows int
	if err := readOnly.QueryRow("select count(*) from orders").Scan(&rows); err != nil {
		t.Errorf("reading as the role: %v", err)
	}

	if _, err := owner.Exec("create table later(id int)"); err != nil {
		t.Fatal(err)
	}
	if err := readOnly.QueryRow("select count(*) from later").Scan(&rows); err != nil {
		t.Errorf("default privileges did not reach a new table: %v", err)
	}

	if got, want := plan(t, owner), plan(t, readOnly); got != want {
		t.Errorf("plans differ under a policy:\nowner\n%s\nrole\n%s", got, want)
	}

	for _, stmt := range []string{
		"insert into orders (id) values (1)",
		"create index orders_probe_idx on orders (status)",
		"explain insert into orders (id) values (1)",
	} {
		if _, err := readOnly.Exec(stmt); err == nil {
			t.Errorf("%s was allowed", stmt)
		}
	}

	// analyze answers a table you do not own with a warning, not an error
	var canSelect, owns bool
	err = readOnly.QueryRow(`select has_table_privilege('orders', 'select'),
		pg_get_userbyid(relowner) = current_user
		from pg_class where relname = 'orders' and relnamespace = 'public'::regnamespace`).
		Scan(&canSelect, &owns)
	if err != nil {
		t.Fatal(err)
	}
	if !canSelect || owns {
		t.Errorf("select %t, owner %t", canSelect, owns)
	}
	if testpg.Major(t) >= 17 {
		var canMaintain bool
		if err := readOnly.QueryRow("select has_table_privilege('orders', 'maintain')").Scan(&canMaintain); err != nil {
			t.Fatal(err)
		}
		if canMaintain {
			t.Error("the role holds maintain")
		}
	}

	if err := db.Drop(ctx); err != nil {
		t.Fatal(err)
	}
	if databaseExists(t, db.Name()) {
		t.Error("database survived its drop")
	}
	if roleExists(t, role.Name()) {
		t.Error("role survived its drop")
	}
}

func plan(t *testing.T, db *sql.DB) string {
	t.Helper()
	rows, err := db.Query("explain select * from orders")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatal(err)
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return strings.Join(lines, "\n")
}

func roleExists(t *testing.T, name string) bool {
	t.Helper()
	var found bool
	err := testpg.Open(t, testpg.Admin(t)).
		QueryRow("select exists (select 1 from pg_roles where rolname = $1)", name).
		Scan(&found)
	if err != nil {
		t.Fatal(err)
	}
	return found
}
