# trac2gitea `import` Packages

These packages implement the data import process.
* `import.data` imports Trac ticket and associated data as Gitea issues
* `import.wiki` imports Trac wiki data into the Gitea project wiki repository

The `import` packages depend on:
* the `accessor` packages for retrieving and storing data
* the `markdown` package for converting Trac markdown to Gitea markdown both for Wiki pages and for issues and their comments

There are no dependencies between the individual `import` packages.