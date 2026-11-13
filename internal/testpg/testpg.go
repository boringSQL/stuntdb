package testpg

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/url"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const envURI = "STUNTDB_TEST_PGURI"

func Admin(t testing.TB) string {
	t.Helper()
	uri := os.Getenv(envURI)
	if uri == "" {
		t.Skipf("%s is not set", envURI)
	}
	return uri
}

func Open(t testing.TB, uri string) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", uri)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func Create(t testing.TB) string {
	t.Helper()
	admin := Admin(t)
	db := Open(t, admin)
	name := scratchName(t)
	if _, err := db.Exec("create database " + name); err != nil {
		t.Fatalf("create database %s: %v", name, err)
	}
	t.Cleanup(func() {
		if _, err := db.Exec("drop database if exists " + name + " with (force)"); err != nil {
			t.Errorf("drop database %s: %v", name, err)
		}
	})
	return databaseURI(t, admin, name)
}

func Major(t testing.TB) int {
	t.Helper()
	var num int
	if err := Open(t, Admin(t)).QueryRow("select current_setting('server_version_num')::int").Scan(&num); err != nil {
		t.Fatalf("server version: %v", err)
	}
	return num / 10000
}

func scratchName(t testing.TB) string {
	t.Helper()
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("scratch name: %v", err)
	}
	return "stuntdb_test_" + hex.EncodeToString(b[:])
}

func databaseURI(t testing.TB, uri, name string) string {
	t.Helper()
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatalf("%s: %v", envURI, err)
	}
	u.Path = "/" + name
	query := u.Query()
	query.Del("dbname")
	u.RawQuery = query.Encode()
	return u.String()
}
