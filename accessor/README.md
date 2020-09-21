# trac2gitea `accessor` Packages

These provide the low-level access primitives:

* `accessor.trac` provides access to Trac data
* `accessor.gitea` provides access to the Gitea project (in particular the database)

There are no dependencies between the individual `accessor` packages.

Each accessor is expressed in terms of an interface `Accessor` with a single, default implementation of that interface `DefaultAccessor`.
This use of interfaces provides a cleaner expression of the accessor functionality and also facilitates testing.
