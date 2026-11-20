package facts_test

import (
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/facts"
)

func TestLabFacts(t *testing.T) {
	f, err := facts.Open("../testdata/lab/facts")
	if err != nil {
		t.Fatal(err)
	}
	if f.Meta.Hash != f.Hash() {
		t.Errorf("meta says %s, files hash to %s", f.Meta.Hash, f.Hash())
	}
	if f.Meta.Masked != "no" {
		t.Errorf("masked = %q", f.Meta.Masked)
	}
	if !strings.Contains(string(f.Schema), "CREATE TABLE public.orders") {
		t.Error("schema.sql is not the lab's")
	}
	if !strings.Contains(string(f.Stats), "\\restrict stuntdb") {
		t.Error("stats.sql kept a random restrict key")
	}
	if len(f.Settings) == 0 {
		t.Error("settings.json is empty")
	}
}
