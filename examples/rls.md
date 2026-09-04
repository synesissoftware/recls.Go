# recls.Go - Example - **rls**


## Summary

A compact Go listing tool inspired by Synesis **`rls`** (`~/.bin/rls.rb`):
recursive search via **recls.Go**, printing modification time, a short
attribute string, size, and path (or path-only in succinct mode).

This is an example, not a full port — multi-root `|` / path-list specs,
directory-size (`-z`), and Windows-only attribute letters are omitted.


## Source

* [examples/rls/main.go](./rls/main.go)


## Execution

```bash
go run ./examples/rls
go run ./examples/rls -s -T . '*.go'
go run ./examples/rls -m --dirs --files -c /tmp
go run ./examples/rls -h
```


<!-- ########################### end of file ########################### -->
