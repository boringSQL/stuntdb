package main

import (
	"context"
	"fmt"
	"io"

	"github.com/boringsql/stuntdb/facts"
)

func runFacts(args []string, stdout, stderr io.Writer) int {
	var source, dir string
	rest := args[1:]
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "-o":
			if i+1 == len(rest) {
				return fail(stderr, "usage: stuntdb facts <uri> -o <dir>")
			}
			dir = rest[i+1]
			i++
		default:
			if source != "" {
				return fail(stderr, "usage: stuntdb facts <uri> -o <dir>")
			}
			source = rest[i]
		}
	}
	if source == "" || dir == "" {
		return fail(stderr, "usage: stuntdb facts <uri> -o <dir>")
	}

	captured, counts, err := facts.Capture(context.Background(), source, dir, "stuntdb "+version)
	if err != nil {
		return fail(stderr, err.Error())
	}

	fmt.Fprintf(stdout, "%-12s%d relations, %d columns\n", "statistics", counts.Relations, counts.Columns)
	settings, err := facts.LoadSettings(captured.Settings)
	if err != nil {
		return fail(stderr, err.Error())
	}
	fmt.Fprintf(stdout, "%-12s%d planner settings\n", "settings", len(settings))
	fmt.Fprintf(stdout, "%-12s%d tables, %d indexes\n", "schema", counts.Tables, counts.Indexes)
	fmt.Fprintf(stdout, "wrote %s   postgres %s   %s   %s\n",
		captured.Dir, captured.Meta.Server, captured.Meta.Captured.Format("2006-01-02"), maskedText(captured.Meta.Masked))
	return exitOK
}

func maskedText(masked string) string {
	switch masked {
	case "yes":
		return "masked"
	case "partly":
		return "partly masked"
	default:
		return "not masked"
	}
}
