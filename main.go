package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	// Embed the zone database so --output-tz works on machines that have
	// no system tzdata (minimal containers, Windows without Go installed).
	_ "time/tzdata"
)

func main() {
	outputTZ := flag.String("output-tz", "UTC", "IANA zone name to print timestamps in, e.g. America/New_York")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: tzfmt [--output-tz ZONE] [file ...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	loc, err := time.LoadLocation(*outputTZ)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tzfmt: invalid --output-tz %q: %v\n", *outputTZ, err)
		os.Exit(2)
	}

	args := flag.Args()
	exitCode := 0

	if len(args) == 0 {
		if !processReader(os.Stdin, "stdin", loc) {
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

		if !processReader(in, label, loc) {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// processReader formats one timestamp per line and reports whether every
// line in this input succeeded, so callers can pick a process exit code.
func processReader(r io.Reader, label string, loc *time.Location) bool {
	ok := true
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		if strings.TrimSpace(raw) == "" {
			continue
		}

		formatted, err := FormatTimestampIn(raw, loc)
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
