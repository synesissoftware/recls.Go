# recls.Go <!-- omit in toc -->

The platform-independent file-system recursive search library, for Go.


[![Language](https://img.shields.io/badge/Language-Go-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![GitHub release](https://img.shields.io/github/v/release/synesissoftware/recls.Go.svg)](https://github.com/synesissoftware/recls.Go/releases/latest)
[![Last Commit](https://img.shields.io/github/last-commit/synesissoftware/recls.Go)](https://github.com/synesissoftware/recls.Go/commits/master)
[![Go](https://github.com/synesissoftware/recls.Go/actions/workflows/go.yml/badge.svg)](https://github.com/synesissoftware/recls.Go/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/synesissoftware/recls.Go.svg)](https://pkg.go.dev/github.com/synesissoftware/recls.Go)


## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Installation](#installation)
- [Components](#components)
  - [Entry](#entry)
  - [Stat](#stat)
  - [Search](#search)
  - [Flags](#flags)
- [Examples](#examples)
- [Project Information](#project-information)
  - [Where to get help](#where-to-get-help)
  - [Contribution guidelines](#contribution-guidelines)
  - [Dependencies](#dependencies)
    - [Development / Testing Dependencies](#development--testing-dependencies)
  - [Related projects](#related-projects)
  - [License](#license)



## Introduction

**recls** — **rec**ursive **ls** — is a platform-independent recursive
file-system search library. **recls.Go** is the Go implementation, shaped
around the recls2 architecture:

* **libpath.Go** — path decomposition (`PathDescriptor`) for Entry fields;
* **stdlib** (`os`, `io/fs`) — directory traversal and metadata;
* **shwild.Go** — shell-compatible pattern matching (not `filepath.Match`).

**recls** is free software released under a BSD-style license.


## Installation

```go
import recls "github.com/synesissoftware/recls.Go"
```


## Components


### Entry

`Entry` embeds **libpath** `PathDescriptor` for path decomposition, adds
search-relative fields, and implements `fs.FileInfo`. Nature predicates are
methods (`Exists()`, `IsDir()`, `IsFile()`, `IsLink()`, `IsHidden()`,
`IsReadonly()`); use `Path()` for the recls name of `FullPath`.


### Stat

```go
e, err := recls.Stat(path, recls.DirectoryParts|recls.MarkDirs)
```

`DetailsLater` allows a path-only entry when the path does not exist.


### Search

Depth-first recursive search. Patterns are matched against the **entry
basename** via **shwild**. `patterns` is a **`PatternSource`**: either a
multi-pattern string (split on `|` or the platform path-list separator —
`:` on Unix, `;` on Windows) or a `[]string` of discrete patterns (not
re-split).

```go
opts := recls.SearchOptions{Flags: recls.Files | recls.Recursive}
for e, err := range recls.Search(root, "*.go", opts) {
    if err != nil { /* ... */ }
    fmt.Println(e.SearchRelativePath)
}
```

Or callback style, with a string or a slice:

```go
err := recls.SearchFunc(root, "*.go|*.md", opts, func(e recls.Entry) error {
    fmt.Println(e.Path())
    return nil
})

err = recls.SearchFunc(root, []string{"*.go", "*.md"}, opts, func(e recls.Entry) error {
    fmt.Println(e.Path())
    return nil
})
```


### Flags

Core flags (aligned with C `RECLS_F_*`): `Files` (default type),
`Directories`, `Links`, `Recursive`, `DirectoryParts`, `MarkDirs`,
`IgnoreHiddenEntries`, `StopOnAccessFailure`, `NoFollowLinks`,
`NoBreakInfiniteLoops`, `DirProgress`, `UseTildeOnNoSearchRoot`,
`DetailsLater`.

Explicitly out of scope for v0.1 (see **TODO.md**): FTP, devices, sockets,
BFS, `RemoveDirectory`, checksums, Windows reparse-dir edge cases.


## Examples

Examples are provided in the `examples` directory, along with a markdown
description for each. A detailed TOC is in [EXAMPLES.md](./EXAMPLES.md).


## Project Information


### Where to get help

[GitHub Page](https://github.com/synesissoftware/recls.Go)


### Contribution guidelines

Defect reports, feature requests, and pull requests are welcome on
https://github.com/synesissoftware/recls.Go.


### Dependencies

* [**libpath.Go**](https://github.com/synesissoftware/libpath.Go/);
* [**shwild.Go**](https://github.com/synesissoftware/shwild.Go/);
* [**ver2go**](https://github.com/synesissoftware/ver2go/);


#### Development / Testing Dependencies

* [**testify**](https://github.com/stretchr/testify);


### Related projects

* [**recls**](https://github.com/synesissoftware/recls/)
* [**recls.NET**](https://github.com/synesissoftware/recls.NET/)
* [**recls.Ruby**](https://github.com/synesissoftware/recls.Ruby/)
* [**recls.Rust**](https://github.com/synesissoftware/recls.Rust/)


### License

**recls.Go** is released under the 3-clause BSD license. See [LICENSE](./LICENSE)
for details.



<!-- ########################### end of file ########################### -->
