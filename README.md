# disk
<img src="docs/img/badges.svg">

The implementation of the [`webtyp.com/files`](https://github.com/webtyp/files) contract on the
machine's disk, for servers and development tools. In the browser, use `webtyp.com/opfs`.

```go
store := disk.Files{}                    // paths as given
data := disk.Files{Root: "/var/lib/app"} // every path inside this directory

var rw files.ReadWriter = store
```

| Operation | Behaviour |
|---|---|
| `ReadFile` | the contents, or `files.ErrNotExist` (never wrapped) when the file is missing |
| `WriteFile` | **atomic**: written to a temporary file in the same directory, then renamed over the target. A reader never sees a truncated file. Mode 0644 |
| `AppendFile` | adds to the end, creating the file (0644) if needed |

It passes `files/conformance`, the same suite every `files` implementation runs.

## Documentation

- [Agent guide](AGENTS.md)
