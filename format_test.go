package main

import (
	"testing"
	"time"
)

func TestClean(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"already clean", "2026-09-21T14:30:00Z", "2026-09-21T14:30:00Z"},
		{"non-breaking space between date and time", "2026-09-21 14:30:00", "2026-09-21 14:30:00"},
		{"unicode minus in offset", "2026-09-21 14:30:00−05:00", "2026-09-21 14:30:00-05:00"},
		{"comma between date and time", "2026-09-21, 14:30:00", "2026-09-21 14:30:00"},
		{"doubled internal whitespace", "2026-09-21   14:30:00", "2026-09-21 14:30:00"},
		{"leading and trailing whitespace", "  2026-09-21 14:30:00  ", "2026-09-21 14:30:00"},
		{"empty string", "", ""},
		{"only whitespace", "   ", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := clean(tc.in)
			if got != tc.want {
				t.Errorf("clean(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"RFC3339Nano", "2026-09-21T14:30:00.123456789Z", "2026-09-21T14:30:00Z"},
		{"RFC3339", "2026-09-21T14:30:00Z", "2026-09-21T14:30:00Z"},
		{"ISO with no zone assumed UTC", "2026-09-21T14:30:00", "2026-09-21T14:30:00Z"},
		{"space separated with Z", "2026-09-21 14:30:00Z", "2026-09-21T14:30:00Z"},
		{"space separated with colon offset", "2026-09-21 09:00:00-05:00", "2026-09-21T14:00:00Z"},
		{"space separated with bare offset", "2026-09-21 09:00:00 -0500", "2026-09-21T14:00:00Z"},
		{"space separated with no zone assumed UTC", "2026-09-21 14:30:00", "2026-09-21T14:30:00Z"},
		{"slash date order", "2026/09/21 14:30:00", "2026-09-21T14:30:00Z"},
		{"US date order", "09/21/2026 14:30:00", "2026-09-21T14:30:00Z"},
		{"day-month-year order", "21-09-2026 14:30:00", "2026-09-21T14:30:00Z"},
		{"non-breaking space and comma cleaned before parse", "2026-09-21, 14:30:00", "2026-09-21T14:30:00Z"},
		{"unicode minus in offset cleaned before parse", "2026-09-21 09:00:00−05:00", "2026-09-21T14:00:00Z"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FormatTimestamp(tc.in)
			if err != nil {
				t.Fatalf("FormatTimestamp(%q) returned error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("FormatTimestamp(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatTimestampIn(t *testing.T) {
	cases := []struct {
		name string
		in   string
		zone string
		want string
	}{
		{"UTC stays Z", "2026-09-21 09:00:00 -0500", "UTC", "2026-09-21T14:00:00Z"},
		{"daylight time zone", "2026-09-21T14:30:00Z", "America/New_York", "2026-09-21T10:30:00-04:00"},
		{"standard time zone", "2026-01-21T14:30:00Z", "America/New_York", "2026-01-21T09:30:00-05:00"},
		{"half hour offset", "2026-09-21T14:30:00Z", "Asia/Kolkata", "2026-09-21T20:00:00+05:30"},
		{"input offset converted", "2026-09-21 09:00:00-05:00", "Europe/Berlin", "2026-09-21T16:00:00+02:00"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			loc, err := time.LoadLocation(tc.zone)
			if err != nil {
				t.Fatalf("LoadLocation(%q): %v", tc.zone, err)
			}
			got, err := FormatTimestampIn(tc.in, loc)
			if err != nil {
				t.Fatalf("FormatTimestampIn(%q) returned error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("FormatTimestampIn(%q, %s) = %q, want %q", tc.in, tc.zone, got, tc.want)
			}
		})
	}
}

func TestFormatTimestampRejectsAmbiguousInput(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"empty input", ""},
		{"whitespace only", "   "},
		{"named zone abbreviation", "2026-09-21 14:30:00 MST"},
		{"garbage", "not a timestamp"},
		{"date only, no time", "2026-09-21"},
		{"two-digit year unsupported", "01/02/03 14:30:00"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FormatTimestamp(tc.in)
			if err == nil {
				t.Errorf("FormatTimestamp(%q) = %q, want error", tc.in, got)
			}
		})
	}
}
