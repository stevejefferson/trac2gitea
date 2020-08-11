# trac2gitea

`trac2gitea` is a command-line tool for migrating [Trac](https://trac.edgewall.org/) projects into [Gitea](https://gitea.io/).

It requires access to the trac project filestore and requires that a corresponding Gitea project has been created into which to migrate the Trac ticket data. If also migrating the associated Trac project wiki then the Gitea project must also have an existing wiki repository (this can however just consist of an empty `Home.md` welcome page).

Trac project data is accessed directly from the Trac database and filestore. While the current code is written for `sqlite` it should be relatively easy to migrate to other SQL databases. 

Gitea project data is also currently accessed directly via the Gitea database (again written for `sqlite` but hopefully not containing much DB-specific SQL). The Gitea access code is highly modularised so hopefully it should be "reasonably easy" to migrate to using the Gitea API in the future. The Gitea wiki is accessed through the Gitea wiki git repository.

## Usage
```
trac2gitea [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo>
Options:
      --db-only                     convert database only
      --default-assignee username   username to assign tickets to when trac assignee is not found in Gitea - defaults to <gitea-user>
      --default-author username     username to attribute content to when trac author is not found in Gitea - defaults to <gitea-user>
      --wiki-dir string             directory into which to checkout (clone) wiki repository - defaults to cwd
      --wiki-only                   convert wiki only
      --wiki-url string             URL of wiki repository - defaults to <server-root-url>/<gitea-user>/<gitea-repo>.wiki.git
```

* `<trac-root>` is the root of the Trac project filestore containing the Trac config file in subdirectory `conf/trac.ini`

* `<gitea-root>` is the root of the Gitea installation

* `<gitea-user>` is the owner of the Gitea project being migrated to

* `<gitea-repo>` is the Gitea repository (project) name being migrated to
