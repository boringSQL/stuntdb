package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFactsUsage(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "no arguments", args: []string{"facts"}},
		{name: "no output", args: []string{"facts", "postgres://localhost/shop"}},
		{name: "two sources", args: []string{"facts", "postgres://localhost/a", "postgres://localhost/b", "-o", "out.facts"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := runFacts(tt.args, &stdout, &stderr); code != exitTrouble {
				t.Errorf("exit code = %d", code)
			}
			if !strings.Contains(stderr.String(), "usage: stuntdb facts") {
				t.Errorf("stderr = %q", stderr.String())
			}
		})
	}
}

func TestFactsExistingDir(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	if code := runFacts([]string{"facts", "postgres://localhost/shop", "-o", dir}, &stdout, &stderr); code != exitTrouble {
		t.Errorf("exit code = %d", code)
	}
	if !strings.Contains(stderr.String(), "already exists") {
		t.Errorf("stderr = %q", stderr.String())
	}
}
