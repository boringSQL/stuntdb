package stuntdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/boringsql/stuntdb/internal/psql"
)

func (d *Database) Load(ctx context.Context, schema []byte) error {
	major, err := d.installExtension(ctx)
	if err != nil {
		return err
	}
	if err := d.checkActive(ctx); err != nil {
		return err
	}
	tool, err := psql.Find()
	if err != nil {
		return err
	}
	if err := tool.Load(ctx, d.uri, schema, major); err != nil {
		return fmt.Errorf("load schema into %s: %w", d.name, err)
	}
	return nil
}

func (d *Database) installExtension(ctx context.Context) (int, error) {
	db, err := d.open(ctx)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var major int
	if err := db.QueryRowContext(ctx, query("server-major")).Scan(&major); err != nil {
		return 0, fmt.Errorf("server version of %s: %w", d.name, err)
	}
	var available bool
	if err := db.QueryRowContext(ctx, query("extension-available")).Scan(&available); err != nil {
		return 0, fmt.Errorf("looking for pg_regresql on %s: %w", d.server.label, err)
	}
	if !available {
		return 0, fmt.Errorf("pg_regresql is not installed on %s: install it for postgres %d", d.server.label, major)
	}
	if _, err := db.ExecContext(ctx, query("create-extension")); err != nil {
		return 0, fmt.Errorf("create extension pg_regresql: %w", err)
	}
	preload := fmt.Sprintf(query("set-preload"), pgx.Identifier{d.name}.Sanitize())
	if _, err := db.ExecContext(ctx, preload); err != nil {
		return 0, fmt.Errorf("preload pg_regresql: %w", err)
	}
	return major, nil
}

func (d *Database) checkActive(ctx context.Context) error {
	db, err := d.open(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	var active *string
	if err := db.QueryRowContext(ctx, query("regresql-active")).Scan(&active); err != nil {
		return fmt.Errorf("checking pg_regresql on %s: %w", d.name, err)
	}
	if active == nil || *active != "on" {
		return fmt.Errorf("pg_regresql is not active on %s: session_preload_libraries did not take", d.name)
	}
	return nil
}
