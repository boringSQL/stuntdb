package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestUpUsage(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "no facts", args: []string{"up"}},
		{name: "no command after --", args: []string{"up", "lab.facts", "--"}},
		{name: "two facts", args: []string{"up", "a.facts", "b.facts"}},
		{name: "schema without a file", args: []string{"up", "lab.facts", "--schema"}},
		{name: "server without a uri", args: []string{"up", "lab.facts", "--server"}},
		{name: "unknown flag", args: []string{"up", "lab.facts", "--set", "work_mem=8MB"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := runUp(tt.args, &stdout, &stderr); code != exitTrouble {
				t.Errorf("exit code = %d", code)
			}
			if !strings.Contains(stderr.String(), "usage: stuntdb up") {
				t.Errorf("stderr = %q", stderr.String())
			}
		})
	}
}

func TestUpMissingFiles(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "facts", args: []string{"up", "../../testdata/lab/missing.facts"}, want: "missing.facts"},
		{name: "schema", args: []string{"up", "../../testdata/lab/facts", "--schema", "../../testdata/lab/missing.sql"}, want: "missing.sql"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := runUp(tt.args, &stdout, &stderr); code != exitTrouble {
				t.Errorf("exit code = %d", code)
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Errorf("stderr = %q, want %q", stderr.String(), tt.want)
			}
		})
	}
}
