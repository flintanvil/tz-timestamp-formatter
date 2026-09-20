# tz-timestamp-formatter

Every system logs timestamps a little differently: `2026-09-21T14:30:00Z`,
`2026-09-21 14:30:00 +0530`, `09/21/2026 14:30:00`, sometimes with a stray
comma or a non-breaking space thrown in by whatever tool exported the file.
This tool reads that mess, one timestamp per line, and prints each one back
out as clean RFC3339 in UTC.

It is not a timezone database. It does not try to guess what "CST" means.
It handles the formats that are unambiguous - ISO 8601 variants and a few
common date orderings, with or without a numeric UTC offset - and rejects
anything it isn't sure about rather than silently misreading it.

## Usage

Build it:

```
go build -o tzfmt .
```

From a file:

```
$ cat sample.txt
2026-09-21T14:30:00Z
2026-09-21 09:00:00 -0500
09/21/2026 14:30:00

$ ./tzfmt sample.txt
2026-09-21T14:30:00Z
2026-09-21T14:00:00Z
2026-09-21T14:30:00Z
```

From stdin, same result:

```
$ cat sample.txt | ./tzfmt
```

Multiple files can be given on the command line, and `-` anywhere in the
list means "read stdin at this point":

```
$ ./tzfmt sample.txt - other.txt < more_timestamps.txt
```

Lines that can't be parsed are reported on stderr with the file name and
line number, and don't stop the rest of the input from being processed. If
any line fails, the process exits with status 1.

Blank lines are skipped. Everything else is treated as one timestamp.

## What counts as UTC

If a line has no offset or `Z` at all (e.g. `2026-09-21 14:30:00`), it's
assumed to already be UTC. That's a real limitation, not a feature - see
the roadmap below.

## Status

Early skeleton. The parser handles a fixed list of layouts; it doesn't yet
take a target output timezone, doesn't handle named zone abbreviations, and
has no test suite yet.
