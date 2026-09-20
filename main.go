package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	exitCode := 0

	if len(args) == 0 {
		if !processReader(os.Stdin, "stdin") {
			exitCode = 1
		}
		os.Exit(exitCode)
	}

	for _, arg := range args {
		var in io.Reader
		var label string

		if arg == "-" {
			in = os.Stdin
			label = "stdin"
		} else {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "tzfmt: %v\n", err)
				exitCode = 1
				continue
			}
			defer f.Close()
			in = f
			label = arg
		}

		if !processReader(in, label) {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// processReader formats one timestamp per line and reports whether every
// line in this input succeeded, so callers can pick a process exit code.
func processReader(r io.Reader, label string) bool {
	ok := true
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		if strings.TrimSpace(raw) == "" {
			continue
		}

		formatted, err := FormatTimestamp(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%d: %v\n", label, lineNum, err)
			ok = false
			continue
		}
		fmt.Println(formatted)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", label, err)
		ok = false
	}

	return ok
}
