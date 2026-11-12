package testpg_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/testpg"
)

func TestCreate(t *testing.T) {
	uri := testpg.Create(t)
	var name string
	if err := testpg.Open(t, uri).QueryRow("select current_database()").Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != databaseName(t, uri) {
		t.Errorf("connected to %s, uri names %s", name, databaseName(t, uri))
	}

	if other := testpg.Create(t); other == uri {
		t.Errorf("two scratch databases share %s", uri)
	}
}

func TestDrop(t *testing.T) {
	var name string
	t.Run("scratch", func(t *testing.T) {
		uri := testpg.Create(t)
		name = databaseName(t, uri)
	})

	var n int
	err := testpg.Open(t, testpg.Admin(t)).
		QueryRow("select count(*) from pg_database where datname = $1", name).
		Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("%s survived its test", name)
	}
}

func TestMajor(t *testing.T) {
	if major := testpg.Major(t); major < 16 {
		t.Errorf("major = %d, want 16 or later", major)
	}
}

func databaseName(t *testing.T, uri string) string {
	t.Helper()
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimPrefix(u.Path, "/")
}
