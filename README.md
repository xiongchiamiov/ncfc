ncfc
====

**ncfc** is a tool for finding directories with large numbers of descendent
files.

I found myself struggling to use `rsync(1)` to copy a directory that had, due
to some poor planning and the design choices of [the build system][bazel],
several million symlinks contained (deep) within it.  This makes rsync rather
unhappy, and produces a lot of overhead with merely the list of updated files
(on initial sync, over 500 MB!), and so I went on a hunt to see where all these
files were coming from.  Even with some convenience wrapper functions, doing a
`find <directory> | wc -l` on every directory in the current level, followed by
moving one layer lower, and repeating ad nauseam, was an unpleasant experience.

Many people go through a similar process with `du(1)` when trying to find large
files until they discover [`ncdu(1)`][ncdu].  Just as ncdu is "**nc**urses
**du**", so ncfc is "**nc**urses **f**ile **c**ount" (or "**nc**urses **f**ile
| w**c**", if you prefer).

[bazel]: http://www.bazel.io/
[ncdu]: https://dev.yorhel.nl/ncdu/man
