# recls.Go - Changes <!-- omit in toc -->


## 0.1.0-alpha1 - 3rd September 2026

* First functional file-system search release (recls2-shaped);
* Added **Entry**, **Stat**, **Search** / **SearchFunc** (DFS, `iter.Seq2`);
* **Search** / **SearchFunc** accept **`PatternSource`** (`string | []string`);
* Path fields from **libpath.Go** `PathDescriptor`; patterns via **shwild.Go**;
* Core search flags: type filter, recursive, hidden, access failure, symlinks,
  infinite-loop guard, mark-dirs, details-later, tilde-on-empty-root;
* Unit and component tests; examples **hard_links**, **libver**, **rls**,
  **search_simple**, **stat**;
* **Entry.LinkCount** (Unix `Stat_t.Nlink`; Windows `BY_HANDLE_FILE_INFORMATION`);
* Migrated version API to **ver2go** 0.2+; CI matrix and lint job;


## 0.0.2 - 2nd September 2026

* Version API migrated toward **ver2go** 0.2 facilities;
* Boilerplate / helper-script updates;


## 0.0.0.4 - 18th August 2025

* GitHub Actions;
* boilerplate;
* documentation;


## 0.0.0.3 - 13th August 2025

* boilerplate;


<!-- ########################### end of file ########################### -->
