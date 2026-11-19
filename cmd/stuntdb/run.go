package main

import (
	"fmt"
	"io"
	"os"
)

const (
	exitOK      = 0
	exitTrouble = 4

	usage = `usage: stuntdb <command> [arguments]

  facts <uri> -o <dir>   statistics, settings and schema from a server
  up <facts>             a throwaway database that plans like production
  version                the build
`
)

type command func(args []string, stdout, stderr io.Writer) int

var commands = map[string]command{
	"version": runVersion,
	"facts":   runFacts,
	"up":      notImplemented,
}

func Run() int {
	return run(os.Args[1:], os.Stdout, os.Stderr)
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return exitTrouble
	}
	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprint(stderr, usage)
		return exitTrouble
	}
	return cmd(args, stdout, stderr)
}

func runVersion(_ []string, stdout, _ io.Writer) int {
	fmt.Fprintln(stdout, build())
	return exitOK
}

func notImplemented(args []string, _ io.Writer, stderr io.Writer) int {
	return fail(stderr, args[0]+": not implemented yet")
}

func fail(stderr io.Writer, msg string) int {
	fmt.Fprintln(stderr, "stuntdb:", msg)
	return exitTrouble
}
