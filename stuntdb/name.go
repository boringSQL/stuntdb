package stuntdb

import (
	"context"
	"fmt"
)

const (
	namePrefix   = "stuntdb_"
	nameAttempts = 10
)

func (s *Server) Name(ctx context.Context) (string, error) {
	for i := 0; i < nameAttempts; i++ {
		name, err := randomName()
		if err != nil {
			return "", err
		}
		var taken bool
		if err := s.db.QueryRowContext(ctx, query("name-taken"), name).Scan(&taken); err != nil {
			return "", fmt.Errorf("checking %s: %w", name, err)
		}
		if !taken {
			return name, nil
		}
	}
	return "", fmt.Errorf("no free database name on %s", s.label)
}

func randomName() (string, error) {
	suffix, err := randomHex(2)
	if err != nil {
		return "", fmt.Errorf("generating a name: %w", err)
	}
	return namePrefix + suffix, nil
}
