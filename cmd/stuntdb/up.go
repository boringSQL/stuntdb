package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/boringsql/stuntdb/facts"
	"github.com/boringsql/stuntdb/stuntdb"
)

var errUpUsage = errors.New("usage: stuntdb up <facts> [--server <uri>] [--schema <file>] [-- <command>]")

type upOptions struct {
	facts   string
	schema  string
	server  string
	command []string
}

func runUp(args []string, stdout, stderr io.Writer) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return upWith(ctx, args, stdout, stderr)
}

func upWith(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	options, err := parseUp(args)
	if err != nil {
		return fail(stderr, err.Error())
	}

	captured, err := facts.Open(options.facts)
	if err != nil {
		return fail(stderr, err.Error())
	}
	settings, err := facts.LoadSettings(captured.Settings)
	if err != nil {
		return fail(stderr, err.Error())
	}
	schema := captured.Schema
	if options.schema != "" {
		if schema, err = os.ReadFile(options.schema); err != nil {
			return fail(stderr, err.Error())
		}
	}

	uri, err := stuntdb.Resolve(options.server)
	if err != nil {
		return fail(stderr, err.Error())
	}
	server, err := stuntdb.Open(ctx, uri)
	if err != nil {
		return fail(stderr, err.Error())
	}
	defer server.Close()

	db, err := server.Create(ctx)
	if err != nil {
		return fail(stderr, err.Error())
	}
	defer drop(db, stderr)

	if err := db.Load(ctx, schema); err != nil {
		return fail(stderr, err.Error())
	}
	statsSkipped, err := db.Inject(ctx, captured.Stats)
	if err != nil {
		return fail(stderr, err.Error())
	}
	settingsSkipped, err := db.Settings(ctx, settings)
	if err != nil {
		return fail(stderr, err.Error())
	}
	role, err := db.ReadOnly(ctx)
	if err != nil {
		return fail(stderr, err.Error())
	}

	fmt.Fprintln(stdout, role.URI())
	fmt.Fprintf(stdout, "plans like production as of %s. Statistics are %s.\n",
		captured.Meta.Captured.Format("2006-01-02"), maskedText(captured.Meta.Masked))
	fmt.Fprintln(stdout, "Empty tables: plan here, do not execute. Read-only. Ctrl-C drops it.")
	for _, entry := range slices.Concat(statsSkipped, settingsSkipped) {
		fmt.Fprintf(stdout, "skipped  %s\n", entry)
	}

	if len(options.command) == 0 {
		<-ctx.Done()
		return exitOK
	}
	return runCommand(options.command, role.URI(), stdout, stderr)
}

func parseUp(args []string) (upOptions, error) {
	var options upOptions
	rest := args[1:]
	for i := 0; i < len(rest); i++ {
		switch arg := rest[i]; {
		case arg == "--":
			options.command = rest[i+1:]
			i = len(rest)
		case arg == "--schema":
			if i+1 == len(rest) {
				return options, errUpUsage
			}
			options.schema = rest[i+1]
			i++
		case arg == "--server":
			if i+1 == len(rest) {
				return options, errUpUsage
			}
			options.server = rest[i+1]
			i++
		case strings.HasPrefix(arg, "-"):
			return options, errUpUsage
		case options.facts != "":
			return options, errUpUsage
		default:
			options.facts = arg
		}
	}
	if options.facts == "" || (options.command != nil && len(options.command) == 0) {
		return options, errUpUsage
	}
	return options, nil
}

// the command runs in stuntdb's process group, so a terminal interrupt reaches
// it first and its own exit code is the one reported
func runCommand(command []string, uri string, stdout, stderr io.Writer) int {
	cmd := exec.Command(command[0], command[1:]...)
	prefix := "DATABASE_URL="
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, prefix) {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, prefix+uri)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err == nil {
		return exitOK
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return fail(stderr, command[0]+": "+err.Error())
	}
	if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		return 128 + int(status.Signal())
	}
	return exitErr.ExitCode()
}

func drop(db *stuntdb.Database, stderr io.Writer) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.Drop(ctx); err != nil {
		fmt.Fprintln(stderr, "stuntdb:", err)
	}
}
