package stuntdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	server *Server
	name   string
}

func (s *Server) Create(ctx context.Context) (*Database, error) {
	name, err := s.Name(ctx)
	if err != nil {
		return nil, err
	}
	stmt := fmt.Sprintf(query("create-database"), pgx.Identifier{name}.Sanitize())
	if _, err := s.db.ExecContext(ctx, stmt); err != nil {
		return nil, fmt.Errorf("create database %s on %s: %w", name, s.label, err)
	}
	return &Database{server: s, name: name}, nil
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Drop(ctx context.Context) error {
	stmt := fmt.Sprintf(query("drop-database"), pgx.Identifier{d.name}.Sanitize())
	if _, err := d.server.db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("drop database %s on %s: %w", d.name, d.server.label, err)
	}
	return nil
}
