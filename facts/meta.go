package facts

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const extendedStatisticsOmitted = "extended statistics: definitions only, pg_dump 18 carries no data for them"

type Counts struct {
	Tables    int
	Columns   int
	Indexes   int
	Relations int
}

const countsQuery = `
select
  (select count(*)
     from pg_class c join pg_namespace n on n.oid = c.relnamespace
    where c.relkind in ('r', 'p', 'm') and n.nspname not in ('pg_catalog', 'information_schema')
      and n.nspname !~ '^pg_'),
  (select count(*)
     from pg_class c join pg_namespace n on n.oid = c.relnamespace
    where c.relkind = 'i' and n.nspname not in ('pg_catalog', 'information_schema')
      and n.nspname !~ '^pg_'),
  (select count(*)
     from pg_attribute a
     join pg_class c on c.oid = a.attrelid
     join pg_namespace n on n.oid = c.relnamespace
    where c.relkind in ('r', 'p', 'm') and a.attnum > 0 and not a.attisdropped
      and n.nspname not in ('pg_catalog', 'information_schema') and n.nspname !~ '^pg_')`

func NewMeta(db *sql.DB, producer, pgDump string) (Meta, error) {
	var server string
	var num int
	err := db.QueryRow("select current_setting('server_version'), current_setting('server_version_num')::int").
		Scan(&server, &num)
	if err != nil {
		return Meta{}, fmt.Errorf("reading server version: %w", err)
	}
	major := num / 10000
	return Meta{
		Format:   Format,
		Server:   strings.Fields(server)[0],
		Major:    major,
		Captured: time.Now().UTC().Truncate(time.Second),
		Masked:   "no",
		Producer: producer,
		PgDump:   pgDump,
		Omitted:  omittedFor(major),
	}, nil
}

func CaptureCounts(db *sql.DB) (Counts, error) {
	var counts Counts
	err := db.QueryRow(countsQuery).
		Scan(&counts.Tables, &counts.Indexes, &counts.Columns)
	if err != nil {
		return Counts{}, fmt.Errorf("counting catalog objects: %w", err)
	}
	counts.Relations = counts.Tables + counts.Indexes
	return counts, nil
}

func omittedFor(major int) []string {
	if major < 19 {
		return []string{extendedStatisticsOmitted}
	}
	return nil
}
