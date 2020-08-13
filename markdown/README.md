# trac2gitea `markdown` Package

This provides the conversion between Trac markdown and Gitea markdown.
The code is heavily based on regular expression matching.

As with the accessors, the markdown converter is expressed in terms of an interface `Converter` with a single, default implementation of that interface `DefaultConverter`.
