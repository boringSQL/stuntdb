package pgdump_test

import (
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/pgdump"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
		err  string
	}{
		{name: "plain", text: "18.6", want: "18.6"},
		{name: "pg_dump output", text: "pg_dump (PostgreSQL) 18.6 (Homebrew)", want: "18.6"},
		{name: "trailing newline", text: "18.6\n", want: "18.6"},
		{name: "no version", text: "pg_dump (PostgreSQL)", err: "no version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := pgdump.ParseVersion(tt.text)
			if tt.err != "" {
				if err == nil || !strings.Contains(err.Error(), tt.err) {
					t.Errorf("error = %v, want %q", err, tt.err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := version.String(); got != tt.want {
				t.Errorf("version = %s, want %s", got, tt.want)
			}
		})
	}
}

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
			version, err := pgdump.ParseVersion(tt.tool)
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
