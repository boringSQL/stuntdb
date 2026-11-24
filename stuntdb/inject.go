package stuntdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (d *Database) Inject(ctx context.Context, stats []byte) ([]Skipped, error) {
	db, err := d.open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := autovacuumOff(ctx, db); err != nil {
		return nil, err
	}
	skipped, err := injectStats(ctx, db, stats)
	if err != nil {
		return nil, err
	}
	missing, err := relationsWithoutStatistics(ctx, db)
	if err != nil {
		return nil, err
	}
	return append(skipped, missing...), nil
}

func autovacuumOff(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, query("relations"))
	if err != nil {
		return fmt.Errorf("listing relations: %w", err)
	}
	type relation struct {
		schema string
		name   string
	}
	var relations []relation
	for rows.Next() {
		var rel relation
		if err := rows.Scan(&rel.schema, &rel.name); err != nil {
			rows.Close()
			return fmt.Errorf("listing relations: %w", err)
		}
		relations = append(relations, rel)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("listing relations: %w", err)
	}

	for _, rel := range relations {
		identifier := pgx.Identifier{rel.schema, rel.name}.Sanitize()
		if _, err := db.ExecContext(ctx, fmt.Sprintf(query("set-autovacuum"), identifier)); err != nil {
			return fmt.Errorf("autovacuum off for %s.%s: %w", rel.schema, rel.name, err)
		}
	}
	return nil
}

func injectStats(ctx context.Context, db *sql.DB, stats []byte) ([]Skipped, error) {
	calls, err := statements(stats)
	if err != nil {
		return nil, err
	}
	var skipped []Skipped
	missing := map[string]bool{}
	for _, call := range calls {
		relation := argValue(call, "schemaname") + "." + argValue(call, "relname")
		subject := relation
		if attribute := argValue(call, "attname"); attribute != "" {
			subject += "." + attribute
		}
		if missing[relation] {
			continue
		}

		var applied bool
		err := db.QueryRowContext(ctx, call).Scan(&applied)
		if err == nil {
			if !applied {
				skipped = append(skipped, Skipped{Subject: subject, Reason: "not applied"})
			}
			continue
		}
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return nil, fmt.Errorf("statistics for %s: %w", subject, err)
		}
		switch pgErr.Code {
		case "42P01":
			missing[relation] = true
			skipped = append(skipped, Skipped{Subject: relation, Reason: "no such relation"})
		case "42703":
			skipped = append(skipped, Skipped{Subject: subject, Reason: "no such column"})
		default:
			return nil, fmt.Errorf("statistics for %s: %w", subject, err)
		}
	}
	return skipped, nil
}

func relationsWithoutStatistics(ctx context.Context, db *sql.DB) ([]Skipped, error) {
	rows, err := db.QueryContext(ctx, query("relations-without-statistics"))
	if err != nil {
		return nil, fmt.Errorf("listing relations without statistics: %w", err)
	}
	defer rows.Close()

	var skipped []Skipped
	for rows.Next() {
		var schema, name string
		if err := rows.Scan(&schema, &name); err != nil {
			return nil, err
		}
		skipped = append(skipped, Skipped{Subject: schema + "." + name, Reason: "no statistics in the facts"})
	}
	return skipped, rows.Err()
}

func statements(stats []byte) ([]string, error) {
	var (
		calls   []string
		current []string
	)
	for _, line := range strings.Split(string(stats), "\n") {
		trimmed := strings.TrimSpace(line)
		if len(current) == 0 {
			if strings.HasPrefix(trimmed, "SELECT * FROM pg_catalog.pg_restore_") {
				current = append(current, line)
			}
			continue
		}
		current = append(current, line)
		if trimmed == ");" {
			calls = append(calls, strings.Join(current, "\n"))
			current = nil
		}
	}
	if len(current) > 0 {
		return nil, errors.New("stats.sql ends inside a statement")
	}
	return calls, nil
}

func argValue(call, name string) string {
	prefix := "'" + name + "', '"
	start := strings.Index(call, prefix)
	if start < 0 {
		return ""
	}
	rest := call[start+len(prefix):]
	end := strings.IndexByte(rest, '\'')
	if end < 0 {
		return ""
	}
	return rest[:end]
}
