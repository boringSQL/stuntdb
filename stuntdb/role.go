package stuntdb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
)

type Role struct {
	name string
	uri  string
}

func (d *Database) ReadOnly(ctx context.Context) (*Role, error) {
	password, err := randomHex(12)
	if err != nil {
		return nil, fmt.Errorf("generating a password: %w", err)
	}
	db, err := d.open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	role := pgx.Identifier{d.name}.Sanitize()
	if _, err := db.ExecContext(ctx, fmt.Sprintf(query("create-role"), role, "'"+password+"'")); err != nil {
		return nil, fmt.Errorf("create role %s on %s: %w", d.name, d.server.label, err)
	}
	d.role = d.name

	if _, err := db.ExecContext(ctx, fmt.Sprintf(query("grant-connect"), pgx.Identifier{d.name}.Sanitize(), role)); err != nil {
		return nil, fmt.Errorf("grant connect to %s: %w", d.name, err)
	}
	schemas, err := userSchemas(ctx, db)
	if err != nil {
		return nil, err
	}
	for _, schema := range schemas {
		identifier := pgx.Identifier{schema}.Sanitize()
		for _, name := range []string{"grant-schema", "grant-select", "grant-default-select"} {
			if _, err := db.ExecContext(ctx, fmt.Sprintf(query(name), identifier, role)); err != nil {
				return nil, fmt.Errorf("%s on %s: %w", name, schema, err)
			}
		}
	}

	uri := d.server.databaseURL(d.name)
	uri.User = url.UserPassword(d.name, password)
	readOnly := &Role{name: d.name, uri: uri.String()}
	if err := readOnly.login(ctx); err != nil {
		return nil, err
	}
	return readOnly, nil
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) URI() string {
	return r.uri
}

func (r *Role) login(ctx context.Context) error {
	db, err := sql.Open("pgx", r.uri)
	if err != nil {
		return fmt.Errorf("open as %s: %w", r.name, err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("log in as %s: %w", r.name, err)
	}
	return nil
}

func userSchemas(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, query("user-schemas"))
	if err != nil {
		return nil, fmt.Errorf("listing schemas: %w", err)
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return schemas, rows.Err()
}
