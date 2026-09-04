# recls.Go - Example - **hard_links**


## Summary

Recursively searches for files and directories under a root (default `.`)
whose hard-link count is greater than 1, printing `linkCount`, a short
kind tag (`file` / `dir` / `link`), and the search-relative path.

On Unix, directories normally report a link count of at least 2 (`.` and
`..`), so most directories will appear; hard-linked **files** (same inode,
multiple names) are the usual signal of interest.


## Source

* [examples/hard_links/main.go](./hard_links/main.go)


## Execution

```bash
go run ./examples/hard_links
go run ./examples/hard_links /tmp
```


<!-- ########################### end of file ########################### -->
