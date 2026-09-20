package main

import (
	"fmt"
	"strings"
	"time"
)

// layouts covers the shapes of timestamp we've actually seen in the wild:
// ISO with and without an offset, space-separated dates, and a couple of
// slash/dash date orderings. Anything with a named zone abbreviation (MST,
// IST, CST...) is deliberately left out - Go will happily "parse" those but
// the offset it picks is not trustworthy without a real tzdata lookup, and
// a wrong silent answer is worse than a rejected line.
var layouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"01/02/2006 15:04:05",
	"02-01-2006 15:04:05",
}

// clean normalises the handful of formatting quirks that show up in
// hand-edited or copy-pasted timestamps without changing the actual value:
// non-breaking spaces, the unicode minus sign some spreadsheets insert,
// stray commas between the date and time, and doubled-up whitespace.
func clean(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case ' ':
			b.WriteRune(' ')
		case '−':
			b.WriteRune('-')
		case ',':
			b.WriteRune(' ')
		default:
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// FormatTimestamp parses a single messy timestamp and returns it as
// RFC3339 in UTC. Inputs with no zone or offset are assumed to already be
// UTC, matching how Go's time.Parse treats a layout with no zone field.
func FormatTimestamp(raw string) (string, error) {
	cleaned := clean(raw)
	if cleaned == "" {
		return "", fmt.Errorf("empty input")
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, cleaned); err == nil {
			return t.UTC().Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("unrecognized timestamp format: %q", raw)
}
