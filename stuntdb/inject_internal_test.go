package stuntdb

import (
	"strings"
	"testing"
)

func TestStatements(t *testing.T) {
	script := `--
-- PostgreSQL database dump
--

\restrict key

SET statement_timeout = 0;
SET transaction_timeout = 0;

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders'
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'status'
);
\unrestrict key
`
	calls, err := statements([]byte(script))
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 {
		t.Fatalf("%d calls", len(calls))
	}
	if !strings.HasPrefix(calls[0], "SELECT * FROM pg_catalog.pg_restore_relation_stats(") || !strings.HasSuffix(calls[0], ");") {
		t.Errorf("first call = %q", calls[0])
	}
	if strings.Contains(calls[1], `\`) || !strings.Contains(calls[1], "'attname', 'status'") {
		t.Errorf("second call = %q", calls[1])
	}
}

func TestStatementsUnterminated(t *testing.T) {
	_, err := statements([]byte("SELECT * FROM pg_catalog.pg_restore_relation_stats(\n\t'schemaname', 'public'\n"))
	if err == nil {
		t.Error("an unterminated statement passed")
	}
}

func TestArgValue(t *testing.T) {
	call := "SELECT * FROM pg_catalog.pg_restore_attribute_stats(\n\t'schemaname', 'public',\n\t'relname', 'orders',\n\t'attname', 'status'\n);"
	if got := argValue(call, "relname"); got != "orders" {
		t.Errorf("relname = %q", got)
	}
	if got := argValue(call, "attname"); got != "status" {
		t.Errorf("attname = %q", got)
	}
	if got := argValue(call, "version"); got != "" {
		t.Errorf("version = %q", got)
	}
}
