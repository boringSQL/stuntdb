package facts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/boringsql/stuntdb/internal/pgdump"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Capture(ctx context.Context, source, dir, producer string) (*Facts, Counts, error) {
	if _, err := os.Stat(dir); err == nil {
		return nil, Counts{}, fmt.Errorf("%s already exists", dir)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, Counts{}, err
	}
	tool, err := pgdump.Find()
	if err != nil {
		return nil, Counts{}, err
	}
	db, err := openReadOnly(ctx, source)
	if err != nil {
		return nil, Counts{}, err
	}
	defer db.Close()

	meta, err := NewMeta(db, producer, tool.Version.String())
	if err != nil {
		return nil, Counts{}, err
	}
	if err := tool.Check(meta.Server); err != nil {
		return nil, Counts{}, err
	}
	settings, err := CaptureSettings(db)
	if err != nil {
		return nil, Counts{}, err
	}
	settingsJSON, err := MarshalSettings(settings)
	if err != nil {
		return nil, Counts{}, err
	}
	counts, err := CaptureCounts(db)
	if err != nil {
		return nil, Counts{}, err
	}
	stats, err := tool.Statistics(ctx, source)
	if err != nil {
		return nil, Counts{}, err
	}
	schema, err := tool.Schema(ctx, source)
	if err != nil {
		return nil, Counts{}, err
	}

	captured := &Facts{Meta: meta, Stats: stats, Settings: settingsJSON, Schema: schema}
	if err := Write(dir, captured); err != nil {
		return nil, Counts{}, err
	}
	return captured, counts, nil
}

func openReadOnly(ctx context.Context, uri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connecting: %w", err)
	}
	if _, err := db.ExecContext(ctx, query("read-only")); err != nil {
		db.Close()
		return nil, fmt.Errorf("setting the session read-only: %w", err)
	}
	return db, nil
}
