package pgdump_test

import (
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/pgdump"
	"github.com/boringsql/stuntdb/internal/pgver"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name   string
		tool   string
		source string
		err    string
	}{
		{name: "same version", tool: "18.6", source: "18.6"},
		{name: "tool ahead", tool: "19.1", source: "18.1"},
		{name: "newer minor server", tool: "18.6", source: "18.7"},
		{name: "too old", tool: "18.5", source: "18.5", err: "too old"},
		{name: "older major", tool: "17.9", source: "18.1", err: "too old"},
		{name: "older than server", tool: "18.6", source: "19.0", err: "cannot dump a 19.0 server"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := pgver.Parse(tt.tool)
			if err != nil {
				t.Fatal(err)
			}
			err = (&pgdump.Tool{Version: version}).Check(tt.source)
			switch {
			case tt.err == "" && err != nil:
				t.Errorf("check = %v", err)
			case tt.err != "" && err == nil:
				t.Errorf("check passed, want %q", tt.err)
			case tt.err != "" && !strings.Contains(err.Error(), tt.err):
				t.Errorf("error = %q, want it to say %q", err, tt.err)
			}
		})
	}
}
