package psql

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/boringsql/stuntdb/internal/pgver"
)

type Tool struct {
	Path    string
	Version pgver.Version
}

func Find() (*Tool, error) {
	path, err := exec.LookPath("psql")
	if err != nil {
		return nil, fmt.Errorf("psql is not on PATH: %w", err)
	}
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil, fmt.Errorf("%s --version: %w", path, err)
	}
	version, err := pgver.Parse(string(out))
	if err != nil {
		return nil, fmt.Errorf("%s --version: %w", path, err)
	}
	return &Tool{Path: path, Version: version}, nil
}

func (t *Tool) Load(ctx context.Context, uri string, script []byte, serverMajor int) error {
	cmd := exec.CommandContext(ctx, t.Path, "-X", "-q", "-v", "ON_ERROR_STOP=1", "--dbname", uri, "-f", "-")
	cmd.Stdin = bytes.NewReader(t.filter(script, serverMajor))
	cmd.Stdout = io.Discard
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func (t *Tool) filter(script []byte, serverMajor int) []byte {
	var out bytes.Buffer
	for _, line := range bytes.SplitAfter(script, []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if serverMajor < 17 && bytes.HasPrefix(trimmed, []byte("SET transaction_timeout")) {
			continue
		}
		if !t.hasRestrict() && (bytes.HasPrefix(trimmed, []byte(`\restrict`)) || bytes.HasPrefix(trimmed, []byte(`\unrestrict`))) {
			continue
		}
		if line, keep := stripAttver(line); keep {
			out.Write(line)
		}
	}
	return out.Bytes()
}

// pg_attver is set by the 2026 minors and refused by servers without them
func stripAttver(line []byte) ([]byte, bool) {
	if !bytes.Contains(line, []byte("pg_attver=")) {
		return line, true
	}
	open := bytes.Index(line, []byte("SET ("))
	end := bytes.LastIndex(line, []byte(");"))
	if open < 0 || end < open {
		return line, true
	}
	var kept []string
	for _, option := range strings.Split(string(line[open+len("SET ("):end]), ",") {
		option = strings.TrimSpace(option)
		if strings.HasPrefix(option, "pg_attver=") {
			continue
		}
		kept = append(kept, option)
	}
	if len(kept) == 0 {
		return nil, false
	}
	var out bytes.Buffer
	out.Write(line[:open])
	out.WriteString("SET (")
	out.WriteString(strings.Join(kept, ", "))
	out.Write(line[end:])
	return out.Bytes(), true
}

// restrict and unrestrict arrived in the 2025 minor releases
func (t *Tool) hasRestrict() bool {
	switch {
	case t.Version.Major >= 18:
		return true
	case t.Version.Major == 17:
		return t.Version.Minor >= 6
	case t.Version.Major == 16:
		return t.Version.Minor >= 10
	}
	return false
}
