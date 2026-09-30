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

## Output zone

Output is UTC by default. `--output-tz` takes an IANA zone name and prints
every timestamp in that zone with its offset at that instant:

```
$ echo 2026-09-21T14:30:00Z | ./tzfmt --output-tz America/New_York
2026-09-21T10:30:00-04:00
```

Flags go before file names. An unknown zone name exits with status 2
before any input is read. The zone database is embedded in the binary, so
this works without system tzdata.

## What counts as UTC

If a line has no offset or `Z` at all (e.g. `2026-09-21 14:30:00`), it's
assumed to already be UTC. That's a real limitation, not a feature - see
the roadmap below.

## Status

Early. The parser handles a fixed list of layouts; it doesn't yet let you
say what zone offset-less input is in, and doesn't handle named zone
abbreviations.
