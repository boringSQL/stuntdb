package stuntdb

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const envServer = "STUNTDB_SERVER"

type Server struct {
	db    *sql.DB
	label string
}

func Resolve(flag string) (string, error) {
	if flag != "" {
		return flag, nil
	}
	if uri := os.Getenv(envServer); uri != "" {
		return uri, nil
	}
	return "", fmt.Errorf("no server: use --server or set %s", envServer)
}

func Open(ctx context.Context, uri string) (*Server, error) {
	cfg, err := pgconn.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	label := serverLabel(cfg)
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to %s: %w", label, err)
	}
	var user string
	var superuser bool
	if err := db.QueryRowContext(ctx, query("current-role")).Scan(&user, &superuser); err != nil {
		db.Close()
		return nil, fmt.Errorf("role on %s: %w", label, err)
	}
	if !superuser {
		db.Close()
		return nil, fmt.Errorf("%s is not a superuser on %s", user, label)
	}
	return &Server{db: db, label: label}, nil
}

func (s *Server) Close() error {
	return s.db.Close()
}

func serverLabel(cfg *pgconn.Config) string {
	if strings.HasPrefix(cfg.Host, "/") {
		return cfg.Host + "/" + cfg.Database
	}
	return net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))) + "/" + cfg.Database
}
