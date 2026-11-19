package facts

import (
	"strings"
	"testing"
)

func TestQueries(t *testing.T) {
	for _, name := range []string{"planner-settings", "catalog-counts", "server-version"} {
		text := query(name)
		if strings.TrimSpace(text) == "" {
			t.Errorf("%s is empty", name)
		}
		if strings.Contains(text, "--") {
			t.Errorf("%s kept its comment: %q", name, text)
		}
	}
	if !strings.Contains(query("planner-settings"), "$1") {
		t.Error("planner-settings lost its parameter")
	}
}
