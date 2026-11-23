package stuntdb_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/testpg"
	"github.com/boringsql/stuntdb/stuntdb"
)

func TestOpen(t *testing.T) {
	open(t)
}

func TestName(t *testing.T) {
	server := open(t)
	name, err := server.Name(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(name, "stuntdb_") || len(name) != len("stuntdb_")+4 {
		t.Fatalf("name = %q", name)
	}

	db := testpg.Open(t, testpg.Admin(t))
	var taken bool
	if err := db.QueryRow("select exists (select 1 from pg_database where datname = $1)", name).Scan(&taken); err != nil {
		t.Fatal(err)
	}
	if taken {
		t.Errorf("%s already exists", name)
	}
}

func TestNonSuperuser(t *testing.T) {
	admin := testpg.Admin(t)
	u, err := url.Parse(admin)
	if err != nil {
		t.Fatal(err)
	}
	db := testpg.Open(t, admin)
	role := scratchRole(t)
	if _, err := db.Exec("create role " + role + " login password 'x'"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := db.Exec("drop role if exists " + role); err != nil {
			t.Errorf("drop role %s: %v", role, err)
		}
	})

	u.User = url.UserPassword(role, "x")
	_, err = stuntdb.Open(context.Background(), u.String())
	if err == nil {
		t.Fatal("connected as a non-superuser")
	}
	if !strings.Contains(err.Error(), "not a superuser") || !strings.Contains(err.Error(), role) {
		t.Errorf("error = %q", err)
	}
	if !strings.Contains(err.Error(), u.Host+u.Path) {
		t.Errorf("error does not name the server: %q", err)
	}
}

func open(t *testing.T) *stuntdb.Server {
	t.Helper()
	server, err := stuntdb.Open(context.Background(), testpg.Admin(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { server.Close() })
	return server
}

func scratchRole(t *testing.T) string {
	t.Helper()
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatal(err)
	}
	return "stuntdb_test_ro_" + hex.EncodeToString(b[:])
}
