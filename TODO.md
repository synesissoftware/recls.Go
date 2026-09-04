# recls.Go - TODO <!-- omit in toc -->


## Table of Contents <!-- omit in toc -->

- [Functional improvements](#functional-improvements)
- [Performance improvements](#performance-improvements)
- [Packaging improvements](#packaging-improvements)



## Functional improvements

* [ ] BFS search order (recls.NET-style);
* [ ] `RemoveDirectory` helper;
* [ ] Devices / sockets type filters;
* [ ] Link count / node index metadata;
* [ ] FTP search (`RECLS_F_PASSIVE_FTP`);
* [ ] CalcChecksum;
* [ ] Windows `AllowReparseDirs` edge cases;
* [ ] Public path helpers (`CombinePaths`, `DeriveRelativePath`, `CanonicalisePath`) — currently internal / stdlib;
* [ ] Upstream `derive_relative_path` / trailing-separator helpers into **libpath.Go** (recls currently carries thin internal copies);
* [ ] Upstream path compare / case-fold helpers into **libpath.Go** (Ruby has `make_compare_path`; use for Windows-aware `pathElementsEqual`);
* [ ] `DirProgress` example and richer progress API;
* [ ] Consider **recls.NET**-style callbacks to filter / decide fate of inaccessible directories (beyond binary StopOnAccessFailure vs skip);
* [x] ~~~Consider generics so `patterns` may be a `string` or an array of strings~~~ - ✅ (`PatternSource`);


## Performance improvements

* [ ] Consider caching compiled patterns across multi-root searches;
* [ ] Optimise patterns: if a wildcards-all (`*`) is present, elide all other patterns;
* [ ] Optional read-ahead / buffered directory reads for large trees;


## Packaging improvements

* [ ] Align **NEWS.md** Details column with umbrella packaging programme;
* [x] ~~~Drop local `replace` directives once published dependency versions cover CI~~~ - ✅;



<!-- ########################### end of file ########################### -->
