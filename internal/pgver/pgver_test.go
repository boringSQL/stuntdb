package pgver_test

import (
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/pgver"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
		err  string
	}{
		{name: "plain", text: "18.6", want: "18.6"},
		{name: "psql output", text: "psql (PostgreSQL) 18.6 (Homebrew)", want: "18.6"},
		{name: "trailing newline", text: "18.6\n", want: "18.6"},
		{name: "no version", text: "psql (PostgreSQL)", err: "no version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := pgver.Parse(tt.text)
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
