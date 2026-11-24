package stuntdb

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (d *Database) Settings(ctx context.Context, settings map[string]string) ([]Skipped, error) {
	names := make([]string, 0, len(settings))
	for name := range settings {
		names = append(names, name)
	}
	sort.Strings(names)

	db, err := d.open(ctx)
	if err != nil {
		return nil, err
	}

	applied := map[string]string{}
	var skipped []Skipped
	for _, name := range names {
		value := settings[name]
		stmt := fmt.Sprintf(query("set-database-setting"),
			pgx.Identifier{d.name}.Sanitize(),
			pgx.Identifier{name}.Sanitize(),
			"'"+strings.ReplaceAll(value, "'", "''")+"'")
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				db.Close()
				return nil, fmt.Errorf("setting %s on %s: %w", name, d.name, err)
			}
			reason := pgErr.Message
			if pgErr.Code == "42704" {
				reason = "unknown to this server"
			}
			skipped = append(skipped, Skipped{Subject: name + "=" + value, Reason: reason})
			continue
		}
		applied[name] = value
	}
	db.Close()

	if len(applied) == 0 {
		return skipped, nil
	}
	notTaken, err := d.readBack(ctx, applied)
	if err != nil {
		return nil, err
	}
	return append(skipped, notTaken...), nil
}

func (d *Database) readBack(ctx context.Context, applied map[string]string) ([]Skipped, error) {
	db, err := d.open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	names := make([]string, 0, len(applied))
	for name := range applied {
		names = append(names, name)
	}
	sort.Strings(names)

	rows, err := db.QueryContext(ctx, query("read-settings"), names)
	if err != nil {
		return nil, fmt.Errorf("reading settings on %s: %w", d.name, err)
	}
	defer rows.Close()

	effective := map[string]string{}
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		effective[name] = value
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var skipped []Skipped
	for _, name := range names {
		if effective[name] != applied[name] {
			skipped = append(skipped, Skipped{Subject: name + "=" + applied[name], Reason: "did not take effect"})
		}
	}
	return skipped, nil
}
