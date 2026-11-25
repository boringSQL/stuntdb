package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		code int
		out  string
		err  string
	}{
		{name: "version", args: []string{"version"}, code: exitOK, out: "stuntdb "},
		{name: "no arguments", code: exitTrouble, err: "usage: stuntdb"},
		{name: "unknown command", args: []string{"nope"}, code: exitTrouble, err: "usage: stuntdb"},
		{name: "facts without arguments", args: []string{"facts"}, code: exitTrouble, err: "stuntdb: usage: stuntdb facts"},
		{name: "up without arguments", args: []string{"up"}, code: exitTrouble, err: "stuntdb: usage: stuntdb up"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != tt.code {
				t.Errorf("exit code = %d, want %d", code, tt.code)
			}
			if !strings.HasPrefix(stdout.String(), tt.out) {
				t.Errorf("stdout = %q, want prefix %q", stdout.String(), tt.out)
			}
			if !strings.HasPrefix(stderr.String(), tt.err) {
				t.Errorf("stderr = %q, want prefix %q", stderr.String(), tt.err)
			}
		})
	}
}
