package stuntdb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	server *Server
	name   string
	uri    string
	role   string
}

func (s *Server) Create(ctx context.Context) (*Database, error) {
	name, err := s.Name(ctx)
	if err != nil {
		return nil, err
	}
	uri := s.databaseURL(name).String()
	stmt := fmt.Sprintf(query("create-database"), pgx.Identifier{name}.Sanitize())
	if _, err := s.db.ExecContext(ctx, stmt); err != nil {
		return nil, fmt.Errorf("create database %s on %s: %w", name, s.label, err)
	}
	return &Database{server: s, name: name, uri: uri}, nil
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Drop(ctx context.Context) error {
	stmt := fmt.Sprintf(query("drop-database"), pgx.Identifier{d.name}.Sanitize())
	if _, err := d.server.db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("drop database %s on %s: %w", d.name, d.server.label, err)
	}
	if d.role == "" {
		return nil
	}
	stmt = fmt.Sprintf(query("drop-role"), pgx.Identifier{d.role}.Sanitize())
	if _, err := d.server.db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("drop role %s on %s: %w", d.role, d.server.label, err)
	}
	return nil
}

func (d *Database) open(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", d.uri)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", d.name, err)
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to %s: %w", d.name, err)
	}
	return db, nil
}

func (s *Server) databaseURL(name string) *url.URL {
	u := *s.url
	u.Path = "/" + name
	params := u.Query()
	params.Del("dbname")
	u.RawQuery = params.Encode()
	return &u
}
