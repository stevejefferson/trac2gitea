# trac2gitea

`trac2gitea` is a command-line utility for migrating [Trac](https://trac.edgewall.org/) projects into [Gitea](https://gitea.io/).

## Scope
At present the following Trac data is converted:
* Trac components, priorities, resolutions, severities, types and versions to Gitea labels
* Trac milestones to Gitea milestones
* Trac tickets to Gitea issues
  * Trac ticket attachments to Gitea issue attachments
  * Trac ticket comments to Gitea issue comments with markdown text conversion
  * Trac ticket labels to Gitea issue labels
* Trac Wiki pages to files in the Gitea wiki repository
  * Markdown text conversion
  * Preservation of Trac wiki page history as git repository commits
* Trac to Gitea markdown conversions (copes with most cases but some Trac constructs may, possibly of necessity, not translate perfectly)
  * link anchors
  * block quotes
  * code blocks (single and multi-line)
  * definition lists
  * Trac bold, italic and underlines to best markdown equivalents
  * headings
  * lists - bulletted, numbered, lettered and roman numbered
  * `[br]` paragraph breaks
  * tables (basic support)
  * Trac links:
    * images
    * `[[url|text]]` style
    * `[url text]` style
    * `http://...` and `https://...` links
    * `htdocs:...` (files are stored in a `htdocs` subdirectory of the Gitea wiki repository)
    * `CamelCase` inter-wiki links
    * `wiki:...` inter-wiki links
    * `attachment:...` current wiki page attachment references (files are stored in a `attachments/<pageName>` subdirectory of the Gitea wiki repository)
    * `attachment:...:wiki:...` wiki attachment references (files are stored in a `attachments/<pageName>` subdirectory of the Gitea wiki repository)
    * `ticket:...` ticket references
    * `comment:...:ticket:...` ticket comment references
    * `attachment:...:ticket:...` ticket attachment references
    * `milestone:...` milestone references
    * `changeset:...` changeset references
    * `source:...` source file references

## Requirements ##
The utility requires access to both the Trac and Gitea filestore.
It retrieves data directly from the Trac database and writes into the Gitea database.
Access to the Gitea project wiki is via by checking out the wiki git repository.

The Gitea project must have been created prior to the migration as must the Gitea project wiki if a Trac wiki is to be converted (this can however just consist of an empty `Home.md` welcome page).

## Usage
```
trac2gitea [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo>
Options:
      --db-only                        convert database only
      --default-assignee username      username to assign tickets to when trac assignee is not found in Gitea - defaults to <gitea-user>
      --default-author username        username to attribute content to when trac author is not found in Gitea - defaults to <gitea-user>
      --default-wiki-author username   username to attribute Wiki content to when trac author is not found in Gitea - defaults to <gitea-user>
      --no-wiki-push                   do not push wiki on completion
      --verbose                        verbose output
      --wiki-convert-predefined        convert Trac predefined wiki pages - by default we skip these
      --wiki-dir string                directory into which to checkout (clone) wiki repository - defaults to cwd
      --wiki-only                      convert wiki only
      --wiki-token string              password/token for accessing wiki repository (ignored if wiki-url provided)
      --wiki-url string                URL of wiki repository - defaults to <server-root-url>/<gitea-user>/<gitea-repo>.wiki.git
```

* `<trac-root>` is the root of the Trac project filestore containing the Trac config file in subdirectory `conf/trac.ini`
* `<gitea-root>` is the root of the Gitea installation
* `<gitea-user>` is the owner of the Gitea project being migrated to
* `<gitea-repo>` is the Gitea repository (project) name being migrated to

## Limitations
The current code is written for `sqlite` (for both the Trac and Gitea databases).

Hopefully, the SQL used by the converter is fairly generic so porting to a different database type should hopefully not be particularly difficult.

For anyone using a different database, the database connections are created in:
  * Trac: `accessor.trac.defaultAccessor.go`, func `CreateDefaultAccessor`
  * Gitea: `accessor.gitea.defaultAccessor.go`, func `CreateDefaultAccessor`

Having changed these, try running the converter and see if any SQL breaks.

All trac database accesses are in package `accessor.trac` and all Gitea database accesses are in package `accessor.gitea`.

## Building
From the root of the source tree run:
```
make
```
This will build the application as an executable `trac2gitea` (in the source tree root directory) and run the tests.

To build the application itself without running the tests, use:
```
make build
```

Missing dependencies can be fetched using:
```
make deps
```

## Acknowledgements
The database migration code is largely derived from [trac2gogs](http://strk.kbt.io/projects/go/trac2gogs/).