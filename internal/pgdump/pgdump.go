package pgdump

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const minimum = "18.6"

var (
	versionPattern = regexp.MustCompile(`(\d+)\.(\d+)`)
	oldest         = Version{Major: 18, Minor: 6}
)

type (
	Tool struct {
		Path    string
		Version Version
	}

	Version struct {
		Major int
		Minor int
	}
)

func Find() (*Tool, error) {
	path, err := exec.LookPath("pg_dump")
	if err != nil {
		return nil, fmt.Errorf("pg_dump is not on PATH: %w", err)
	}
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil, fmt.Errorf("%s --version: %w", path, err)
	}
	version, err := ParseVersion(string(out))
	if err != nil {
		return nil, fmt.Errorf("%s --version: %w", path, err)
	}
	return &Tool{Path: path, Version: version}, nil
}

func (t *Tool) Statistics(ctx context.Context, uri string) ([]byte, error) {
	return t.dump(ctx, uri, "--statistics-only")
}

func (t *Tool) Schema(ctx context.Context, uri string) ([]byte, error) {
	return t.dump(ctx, uri, "--schema-only", "--no-owner", "--no-acl")
}

func (t *Tool) Check(source string) error {
	if t.Version.less(oldest) {
		return fmt.Errorf("pg_dump %s is too old: statistics need %s or later", t.Version, minimum)
	}
	server, err := ParseVersion(source)
	if err != nil {
		return fmt.Errorf("server version %q: %w", source, err)
	}
	if t.Version.Major < server.Major {
		return fmt.Errorf("pg_dump %s cannot dump a %s server", t.Version, source)
	}
	return nil
}

func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor)
}

func (v Version) less(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	return v.Minor < other.Minor
}

func ParseVersion(text string) (Version, error) {
	match := versionPattern.FindStringSubmatch(text)
	if match == nil {
		return Version{}, fmt.Errorf("no version in %q", strings.TrimSpace(text))
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return Version{}, err
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor}, nil
}

func (t *Tool) dump(ctx context.Context, uri string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, t.Path, append(args, "--dbname", uri)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pg_dump %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}
