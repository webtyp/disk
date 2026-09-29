# Agent Guide — `webtyp/disk`

## What this library is

The disk implementation of `webtyp.com/files`. It is **backend code**: it legitimately uses
`os` and `path/filepath`. Do not "fix" those imports, and do not import it from browser code
(use `webtyp.com/opfs` there).

## Rules

- `WriteFile` stays atomic (temporary file in the same directory, then rename). A plain
  `os.WriteFile` truncates first and writes second, and a reader in that window cannot tell an
  emptied file from an empty one. Real projects lost their `.env` keys that way.
- A missing file is `files.ErrNotExist`, unwrapped.
- Every change keeps `conformance.Run` green.

## The builds that define "done"

```bash
go vet ./...
gotest
```
