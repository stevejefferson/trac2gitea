# trac2gitea `import` Packages

These implement the data import process.
* `import.issue` imports Trac ticket and associated data as Gitea issues
* `import.wiki` imports Trac wiki data into the Gitea project wiki repository

The `import` packages depend on the `accessor` packages for retrieving and storing data.
They also depend on the `markdown` package for converting Trac markdown to Gitea markdown both for Wiki pages and for issues and their comments.

There are no dependencies between the separate `import` packages.