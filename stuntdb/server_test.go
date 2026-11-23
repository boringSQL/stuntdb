package stuntdb

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	t.Run("flag", func(t *testing.T) {
		t.Setenv(envServer, "postgres://env/shop")
		uri, err := Resolve("postgres://flag/shop")
		if err != nil || uri != "postgres://flag/shop" {
			t.Errorf("uri = %q, err = %v", uri, err)
		}
	})
	t.Run("environment", func(t *testing.T) {
		t.Setenv(envServer, "postgres://env/shop")
		uri, err := Resolve("")
		if err != nil || uri != "postgres://env/shop" {
			t.Errorf("uri = %q, err = %v", uri, err)
		}
	})
	t.Run("database url", func(t *testing.T) {
		t.Setenv(envServer, "")
		t.Setenv("DATABASE_URL", "postgres://app/shop")
		uri, err := Resolve("")
		if err == nil {
			t.Fatalf("resolved %q from DATABASE_URL", uri)
		}
		if !strings.Contains(err.Error(), "--server") || !strings.Contains(err.Error(), envServer) {
			t.Errorf("error = %q", err)
		}
	})
}

func TestRandomName(t *testing.T) {
	for i := 0; i < 100; i++ {
		name, err := randomName()
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(name, namePrefix) || len(name) != len(namePrefix)+4 {
			t.Fatalf("name = %q", name)
		}
		if _, err := hex.DecodeString(name[len(namePrefix):]); err != nil {
			t.Fatalf("name = %q: %v", name, err)
		}
	}
}

func TestOpenRejectsDSN(t *testing.T) {
	_, err := Open(context.Background(), "host=localhost dbname=shop")
	if err == nil || !strings.Contains(err.Error(), "not a uri") {
		t.Errorf("error = %v", err)
	}
}
