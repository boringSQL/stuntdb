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

func NewMeta(db *sql.DB, producer, pgDump string) (Meta, error) {
	var server string
	var num int
	err := db.QueryRow(query("server-version")).Scan(&server, &num)
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
	err := db.QueryRow(query("catalog-counts")).
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
