package psql

import (
	"slices"
	"strings"
	"testing"

	"github.com/boringsql/stuntdb/internal/pgver"
)

const script = `\restrict stuntdb
SET statement_timeout = 0;
SET transaction_timeout = 0;
create table rides(id bigint primary key);
\unrestrict stuntdb
`

func TestFilterAttver(t *testing.T) {
	attver := []byte("ALTER TABLE ONLY public.t ALTER COLUMN id SET (pg_attver=127566736);\n" +
		"ALTER TABLE ONLY public.t ALTER COLUMN code SET (pg_attver=127566736, n_distinct=5);\n" +
		"ALTER TABLE ONLY public.t ALTER COLUMN note SET (n_distinct=7);\n")
	got := (&Tool{Version: pgver.Version{Major: 18}}).filter(attver, 18)
	want := "ALTER TABLE ONLY public.t ALTER COLUMN code SET (n_distinct=5);\n" +
		"ALTER TABLE ONLY public.t ALTER COLUMN note SET (n_distinct=7);\n"
	if string(got) != want {
		t.Errorf("filtered\n%s\nwant\n%s", got, want)
	}
}

func TestFilter(t *testing.T) {
	all := []string{"\\restrict stuntdb", "SET statement_timeout = 0;", "SET transaction_timeout = 0;", "create table rides(id bigint primary key);", "\\unrestrict stuntdb"}
	noTimeout := []string{"\\restrict stuntdb", "SET statement_timeout = 0;", "create table rides(id bigint primary key);", "\\unrestrict stuntdb"}
	noRestrict := []string{"SET statement_timeout = 0;", "SET transaction_timeout = 0;", "create table rides(id bigint primary key);"}
	statementsOnly := []string{"SET statement_timeout = 0;", "create table rides(id bigint primary key);"}

	tests := []struct {
		name   string
		client pgver.Version
		server int
		want   []string
	}{
		{name: "18 server and client", client: pgver.Version{Major: 18, Minor: 6}, server: 18, want: all},
		{name: "16 server drops the timeout", client: pgver.Version{Major: 18, Minor: 6}, server: 16, want: noTimeout},
		{name: "17 server keeps the timeout", client: pgver.Version{Major: 18, Minor: 6}, server: 17, want: all},
		{name: "16.4 client drops restrict", client: pgver.Version{Major: 16, Minor: 4}, server: 18, want: noRestrict},
		{name: "16.10 client keeps restrict", client: pgver.Version{Major: 16, Minor: 10}, server: 18, want: all},
		{name: "17.5 client drops restrict", client: pgver.Version{Major: 17, Minor: 5}, server: 18, want: noRestrict},
		{name: "17.6 client keeps restrict", client: pgver.Version{Major: 17, Minor: 6}, server: 18, want: all},
		{name: "15 client drops both", client: pgver.Version{Major: 15, Minor: 8}, server: 16, want: statementsOnly},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := &Tool{Version: tt.client}
			got := strings.Split(strings.TrimSpace(string(tool.filter([]byte(script), tt.server))), "\n")
			if !slices.Equal(got, tt.want) {
				t.Errorf("kept\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(tt.want, "\n"))
			}
		})
	}
}
