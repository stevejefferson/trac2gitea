// Copyright 2016 by Sandro Santilli <strk@kbt.io>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA
// 02110-1301  USA

package main

/*
 * Globals
 */
import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-ini/ini"
	"github.com/spf13/pflag"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var defaultAuthorID int64
var defaultAssigneeID int64
var tracDB, giteaDB *sql.DB
var giteaRepoID int64

var err error

// + + + +

/*
 * Utility routines
 */
func setLogFormat() {
	log.SetFlags(log.Ldate)
	log.SetFlags(log.Ltime)
	log.SetFlags(log.Lshortfile)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}

	return cerr
}

// + + + +

/*
 * Config file access
 */
func fetchConfig(configPath string) *ini.File {
	stat, err := os.Stat(configPath)
	if err != nil {
		return nil
	}
	if stat.IsDir() == true {
		return nil
	}

	config, err := ini.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func getConfigValue(config *ini.File, sectionName string, configName string) string {
	if config == nil {
		return ""
	}

	configValue, err := config.Section(sectionName).GetKey(configName)
	if err != nil {
		return ""
	}

	return configValue.String()
}
// + + + +

/*
 * Arguments
 */
var defaultAssignee string
var defaultAuthor string
var dbOnly bool
var wikiOnly bool
var tracRootDir string
var giteaRootDir string
var giteaUserName string
var giteaRepoName string
var giteaWikiRepoDir string

func parseArgs() {
	defaultAssigneeParam := pflag.String("default-assignee", "",
		"`username` to assign tickets to when trac assignee is not found in Gitea")
	defaultAuthorParam := pflag.String("default-author", "",
		"`username` to attribute content to when trac author is not found in Gitea")
	dbOnlyParam := pflag.Bool("db-only", false,
		"convert database only")
	wikiOnlyParam := pflag.Bool("wiki-only", false,
		"convert wiki only")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <trac_root> <gitea_root> <gitea_user> <gitea_repo_name> <gitea_wiki_repo_dir>\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	dbOnly = *dbOnlyParam;
	wikiOnly = *wikiOnlyParam;
	if dbOnly && wikiOnly {
		log.Fatal("Cannot generate only database AND only wiki!")
	}

	if pflag.NArg() < 5 {
		pflag.Usage()
		os.Exit(1)
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUserName = pflag.Arg(2)
	giteaRepoName = pflag.Arg(3)
	giteaWikiRepoDir = pflag.Arg(4)

	defaultAssignee = *defaultAssigneeParam;
	defaultAuthor = *defaultAuthorParam;
}

func validateArgs() {
	stat, err := os.Stat(tracRootDir)
	if err != nil {
		log.Fatal(err)
	}
	if stat.IsDir() != true {
		log.Fatal("Trac root directory is not a directory: ", tracRootDir)
	}
	fetchTracConfig(tracRootDir)

	stat, err = os.Stat(giteaRootDir)
	if err != nil {
		log.Fatal(err)
	}
	if stat.IsDir() != true {
		log.Fatal("Gitea root path is not a directory: ", giteaRootDir)
	}
	giteaFetchConfig(giteaRootDir)

	if !dbOnly {
		stat, err = os.Stat(giteaWikiRepoDir);
		if err != nil {
			log.Fatal(err)
		}
		if stat.IsDir() != true {
			log.Fatal("Gitea wiki repo directory is not a directory: ", giteaWikiRepoDir)
		}
	}
}

// + + + +

/*
 * Gitea Operations
 */
func getGiteaDB() *sql.DB {
	// extract path to gitea DB - currently sqlite-specific...
	giteaDatabasePath := getGiteaConfigValue("database", "PATH")
	giteaDatabase, err := sql.Open("sqlite3", giteaDatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using Gitea database %s\n", giteaDatabasePath);
	return giteaDatabase
}

var giteaMainConfig *ini.File
var giteaCustomConfig *ini.File
func giteaFetchConfig(rootDir string) {
	if giteaMainConfig != nil || giteaCustomConfig  != nil {
		return		// config already fetched
	}

	giteaMainConfigPath := "/etc/gitea/conf/app.ini"
	giteaMainConfig = fetchConfig(giteaMainConfigPath);
	giteaCustomConfigPath := fmt.Sprintf("%s/custom/conf/app.ini", rootDir)
	giteaCustomConfig = fetchConfig(giteaCustomConfigPath);
	if giteaMainConfig == nil && giteaCustomConfig == nil {
		log.Fatal("cannot find Gitea config in " + giteaMainConfigPath + " or " + giteaCustomConfigPath)
	}
}

func getGiteaConfigValue(sectionName string, configName string) string {
	configValue := getConfigValue(giteaCustomConfig, sectionName, configName)
	if configValue == "" {
		configValue = getConfigValue(giteaMainConfig, sectionName, configName)
	}

	return configValue
}

func giteaFindUserID(nameOrAddress string) int64 {
	if strings.Trim(nameOrAddress, " ") == "" {
		return -1
	}

	var id int64
	err := giteaDB.QueryRow(`SELECT id FROM user WHERE lower_name = $1 or email = $1`, nameOrAddress).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

func giteaFindAdminUserName() string {
	row := giteaDB.QueryRow(`
		SELECT lower_name FROM user WHERE is_admin ORDER BY id LIMIT 1;
		`)

	var name string
	err := row.Scan(&name)
	if err != nil {
		log.Fatal("No admin user found in Gitea")
	}

	return name
}

func giteaFindUserOrAdminID(userName string, adminUserID int64) int64 {
	userID := adminUserID
	if userName != "" {
		userID = giteaFindUserID(userName)
		if userID == -1 {
			log.Fatal("Cannot find gitea user ", userName)
		}
	}

	return userID
}

func giteaFindRepoID(userName string, repoName string) int64 {
	row := giteaDB.QueryRow(`
		SELECT r.id FROM repository r, user u WHERE r.owner_id =
			u.id AND u.name = $1 AND r.name = $2
		`, userName, repoName)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Fatal("No Gitea repository " + repoName + " found for user " + userName);
	}

	return id
}

func giteaUpdateRepoIssueCount(count int, closedCount int) {
	// Update issue count for repo
	if count > 0 {
		_, err = giteaDB.Exec(`
			UPDATE repository SET num_issues = num_issues+$1
				WHERE id = $2`,
			count, giteaRepoID)
		if err != nil {
			log.Fatal(err)
		}
	}
	if closedCount > 0 {
		_, err = giteaDB.Exec(`
			UPDATE repository
				SET num_closed_issues = num_closed_issues+$1
				WHERE id = $2`,
			closedCount, giteaRepoID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func giteaAddLabel(label string, color string) {
	_, err := giteaDB.Exec(`
		INSERT INTO label(repo_id,name,color)
			SELECT $1,$2, $3 WHERE
			NOT EXISTS ( SELECT * FROM label WHERE repo_id = $1 AND name = $2 )`,
		giteaRepoID, label, color)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added label", label)
}

func giteaAddIssueLabel(issueID int64, label string) {
	_, err := giteaDB.Exec(`
		INSERT INTO issue_label(issue_id, label_id)
			SELECT $1, (SELECT id FROM label where repo_id = $2 and name = $3)`,
		issueID, giteaRepoID, label)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added label %s for issue %d\n", label, issueID);
}

func giteaAddMilestone(name string, content string, closed bool, deadlineTS int64, closedTS int64) {
	_, err := giteaDB.Exec(`
		INSERT INTO
			milestone(repo_id,name,content,is_closed,deadline_unix,closed_date_unix)
			SELECT $1,$2,$3,$4,$5,$6 WHERE
				NOT EXISTS (SELECT * FROM milestone WHERE repo_id = $1 AND name = $2)`,
		giteaRepoID, name, content, closed, deadlineTS, closedTS)
	if err != nil {
		log.Fatal(err)
	}
}

func giteaAddIssue(
		ticketID int64, summary string, reporterID int64,
		milestone string, ownerID sql.NullString, owner string,
		closed bool, description string, created int64) int64 {
	_, err := giteaDB.Exec(`
		INSERT INTO issue('index', repo_id, name, poster_id, milestone_id, original_author_id, original_author, is_pull, is_closed, content, created_unix)
			SELECT $1, $2, $3, $4, (SELECT id FROM milestone WHERE repo_id = $2 AND name = $5), $6, $7, false, $8, $9, $10`,
		ticketID, giteaRepoID, summary, reporterID, milestone, ownerID, owner, closed, description, created)
	if err != nil {
		log.Fatal(err)
	}

	var gid int64
	err = giteaDB.QueryRow(`SELECT last_insert_rowid()`).Scan(&gid)
	if err != nil {
		log.Fatal(err)
	}

	return gid
}

func giteaSetIssueUpdateTime(issueID int64, updateTime int64) {
	_, err = giteaDB.Exec(`UPDATE issue SET updated_unix = MAX(updated_unix,$1) WHERE id = $2`, updateTime, issueID)
	if err != nil {
		log.Fatal(err)
	}
}

func giteaAddComment(issueID int64, authorID int64, comment string, time int64) int64 {
	_, err := giteaDB.Exec(`
		INSERT INTO comment(
			type, issue_id, poster_id, content, created_unix, updated_unix)
			VALUES ( 0, $1, $2, $3, $4, $4 )`,
		issueID, authorID, comment, time)
	if err != nil {
		log.Fatal(err)
	}

	var commentID int64
	err = giteaDB.QueryRow(`SELECT last_insert_rowid()`).Scan(&commentID)
	if err != nil {
		log.Fatal(err)
	}

	return commentID
}

func giteaAddAttachment(uuid string, issueID int64, commentID int64, fname string, time int64) {
	_, err = giteaDB.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, uuid, issueID, commentID, fname, time)
	if err != nil {
		log.Fatal(err)
	}
}

func giteaAttachmentURL(uuid string) string {
	baseURL := getGiteaConfigValue("server", "ROOT_URL")
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}

func giteaAttachmentPath(uuid string) string {
	d1 := uuid[0:1]
	d2 := uuid[1:2]
	// TODO: seek for PATH under [attachment]
	//       in giteaRootPath/custom/conf/app.ini
	subpath := "data/attachments"
	return fmt.Sprintf("%s/%s/%s/%s/%s", giteaRootDir, subpath, d1, d2, uuid)
}

func giteaWriteWikiPageVersion(pageName string, markdownText string, version int64, comment string, time int64) {
	pagePath := filepath.Join(giteaWikiRepoDir, pageName + ".md")
	file, err := os.Create(pagePath)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(markdownText)
}

func giteaCopyFile(externalFilePath string, giteaPath string) {
	_, err = os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: cannot copy non-existant file: \"%s\"\n", externalFilePath)
		return
	}

	err = os.MkdirAll(path.Dir(giteaPath), 0775)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(giteaPath)
	if os.IsExist(err) {
		return	// if gitea path exists we'll just assume we've already created it as part of this run
	}

	err = copyFile(externalFilePath, giteaPath)
	if err != nil {
		log.Fatal(err)
	}
}

// + + + +

/*
 * Trac operations
 */
var tracConfig *ini.File

func fetchTracConfig(rootDir string) {
	if tracConfig != nil {
		return		// config already fetched
	}

	tracIniPath := fmt.Sprintf("%s/conf/trac.ini", rootDir)
	tracConfig = fetchConfig(tracIniPath)
	if tracConfig == nil {
		log.Fatal("cannot find trac ini file in " + tracIniPath)
	}
}

func getTracConfigValue(sectionName string, configName string) string {
	return getConfigValue(tracConfig, sectionName, configName)
}

func getTracDB() *sql.DB {
	// extract path to trac DB - currently sqlite-specific...
	tracDatabaseString := getTracConfigValue("trac", "database")
	tracDatabaseSegments := strings.SplitN(tracDatabaseString, ":", 2)
	tracDatabasePath := tracDatabaseSegments[1];
	if !filepath.IsAbs(tracDatabasePath) {
		tracDatabasePath = filepath.Join(tracRootDir, tracDatabasePath)
	}

	tracDatabase, err := sql.Open("sqlite3", tracDatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using Trac database %s\n", tracDatabasePath);
	return tracDatabase
}

func encodeSha1(str string) string {
	// Encode string to sha1 hex value.
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func tracAttachmentPath(tid int64, fname string) string {
	ticketDir := encodeSha1(fmt.Sprintf("%d", tid))
	ticketSub := ticketDir[0:3]

	pathFile := encodeSha1(fname)
	pathExt := path.Ext(fname)

	return fmt.Sprintf("%s/attachments/ticket/%s/%s/%s%s", tracRootDir, ticketSub, ticketDir, pathFile, pathExt)
}

// + + + +

/*
 * Trac to markdown conversion
 */
// @param prefix is the prefix to a link, like "ticket:5" for
//        a link like "ticket:5:inside_the_nest.png"
func resolveTracLink(link, prefix string) string {
	// 'http...' links are left as-is
	if strings.HasPrefix(link, "http") {
		return link
	}

	// 'htdocs:...' links refer to trac htdocs directory
	if strings.HasPrefix(link, "htdocs:") {
		htdocsPath := strings.Replace(link, "htdocs:", "", -1)
		tracHtdocsPath := filepath.Join(tracRootDir, "htdocs", htdocsPath)
		giteaHtdocsPath := filepath.Join(giteaWikiRepoDir, "htdocs", htdocsPath)
		giteaCopyFile(tracHtdocsPath, giteaHtdocsPath)
		return htdocsPath
	}

	// 'ticket:...' links refer to attachments
	if strings.HasPrefix(prefix, "ticket:") {
		tid := strings.Replace(prefix, "ticket:", "", -1)

		// Find issue id
		var gid int64
		err := giteaDB.QueryRow(`
select id from issue where repo_id = $1 and index = $2
			`, giteaRepoID, tid).Scan(&gid)
		if err != nil {
			log.Fatal(err)
		}

		// Find attachment
		var uuid string
		err = giteaDB.QueryRow(`
select uuid from attachment where issue_id = $1 and name = $2
			`, gid, link).Scan(&uuid)
		if err != nil {
			log.Fatal(err)
		}

		return giteaAttachmentURL(uuid)
	}

	// TODO 'wiki:...' links

	fmt.Fprintf(os.Stderr, "WARNING: cannot resolve trac link %s with prefix '%s'", link, prefix)
	return link
}

func tracCodeBlockToMarkdown(in, linkprefix string) string {
	// convert single line {{{...}}} to `...`
	out := in
	re := regexp.MustCompile("{{{([^\n]+?)}}}")
	out = re.ReplaceAllString(out, "`$1`")

	// convert multi-line {{{...}}} to tab-indented lines
	re = regexp.MustCompile("(?s){{{(.+?)}}}")
	out = re.ReplaceAllStringFunc(out, func(m string) string {
		lines := strings.Split(m, "\n")
		for i := range lines {
			l := lines[i]
			l = strings.Replace(l, "{{{", "", -1)
			l = strings.Replace(l, "}}}", "", -1)
			l = "\t" + l
			lines[i] = l
		}
		return strings.Join(lines, "\n")
	})

	return out
}

func tracImageReferenceToMarkdown(in, linkprefix string) string {
	regex := regexp.MustCompile(`\[\[Image\(([^,\)]*)[^\)]*\)\]\]`)
	out := regex.ReplaceAllStringFunc(in, func(m string) string {
		u := regex.ReplaceAllString(m, "$1")
		u = resolveTracLink(u, linkprefix)
		return fmt.Sprintf("![](%s)", u)
	})

	return out
}

func tracHeadingToMarkdown(in string, tracDelimiter string, markdownDelimiter string) string {
	if len(tracDelimiter) >= 6 {
		return in
	}

	tracDelimiter = tracDelimiter + "#"
	markdownDelimiter = markdownDelimiter + "="
	regexStr := tracDelimiter + `\s+(\S+)\s+A` + tracDelimiter
	regex := regexp.MustCompile(regexStr)
	out := regex.ReplaceAllStringFunc(in, func(m string) string {
		return ""
	})

	return out
}

func tracHeadingsToMarkdown(in string) string {
	return tracHeadingToMarkdown(in, "", "")
}

func tracTextToMarkdown(in, linkPrefix string) string {
	out := tracCodeBlockToMarkdown(in, linkPrefix)
	out = tracImageReferenceToMarkdown(out, linkPrefix)
	out = tracHeadingsToMarkdown(out)

	// TODO:
	//    body.gsub!(/\=\=\=\=\s(.+?)\s\=\=\=\=/, '### \1')
	//    body.gsub!(/\=\=\=\s(.+?)\s\=\=\=/, '## \1')
	//    body.gsub!(/\=\=\s(.+?)\s\=\=/, '# \1')
	//    body.gsub!(/\=\s(.+?)\s\=[\s\n]*/, '')
	//    body.gsub!(/\[(http[^\s\[\]]+)\s([^\[\]]+)\]/, '[\2](\1)')
	//    body.gsub!(/\!(([A-Z][a-z0-9]+){2,})/, '\1')
	//    body.gsub!(/'''(.+)'''/, '*\1*')
	//    body.gsub!(/''(.+)''/, '_\1_')
	//    body.gsub!(/^\s\*/, '*')
	//    body.gsub!(/^\s\d\./, '1.')

	return out
}

// + + + +

/*
 * Wiki import
 */
func importWiki() {
	rows, err := tracDB.Query(`
		SELECT w1.name, w1.text, w1.comment, w1.version, w1.time
			FROM wiki w1
			WHERE w1.version = (SELECT MAX(w2.version) FROM wiki w2 WHERE w1.name = w2.name)`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var name string
		var text string
		var commentStr sql.NullString
		var version int64
		var time int64
		if err := rows.Scan(&name, &text, &commentStr, &version, &time); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Converting Wiki page %s, version %d\n", name, version)
		comment := ""
		if !commentStr.Valid {
			comment = commentStr.String
		}
		markdownText := tracTextToMarkdown(text, "")
		giteaWriteWikiPageVersion(name, markdownText, version, comment, time)
	}
}

// + + + +

/*
 * Issue Import
 */
func importLabels(tracQuery string, labelPrefix string, labelColor string) {
	rows, err := tracDB.Query(tracQuery)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			log.Fatal(err)
		}
		lbl := labelPrefix + " / " + val
		giteaAddLabel(lbl, labelColor)
	}
}

func importComponents() {
	importLabels(`SELECT name FROM component`, "Component", "#fbca04");
}

func importPriorities() {
	importLabels(`SELECT DISTINCT priority FROM ticket`, "Priority", "#207de5");
}

func importSeverities() {
	importLabels(`SELECT DISTINCT COALESCE(severity,'') FROM ticket`, "Severity", "#eb6420");
}

func importVersions() {
	importLabels(`SELECT DISTINCT COALESCE(version,'') FROM ticket UNION
                        SELECT COALESCE(name,'') FROM version`, "Version", "#009800");
}

func importTypes() {
	importLabels(`SELECT DISTINCT type FROM ticket`, "Type", "#e11d21");
}

func importResolutions() {
	importLabels(`SELECT DISTINCT resolution FROM ticket WHERE trim(resolution) != ''`, "Resolution", "#9e9e9e");
}

func importMilestones() {
	// NOTE: trac timestamps are to the microseconds, we just need seconds
	rows, err := tracDB.Query(`
		SELECT COALESCE(name,''), CAST(due*1e-6 AS int8), CAST(completed*1e-6 AS int8), description
			FROM milestone UNION
			SELECT distinct(COALESCE(milestone,'')),0,0,''
				FROM ticket
				WHERE COALESCE(milestone,'') NOT IN ( select COALESCE(name,'') from milestone )`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var completed, due int64
		var nam, desc string
		if err := rows.Scan(&nam, &due, &completed, &desc); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Adding milestone", nam)
		giteaAddMilestone(nam, desc, completed != 0, due, completed)
	}
}

func importTicket(ticketID int64, created int64, owner string, reporter string, milestone string, closed bool, summary string, description string) int64 {
	description = tracTextToMarkdown(description, "")	// in Trac format

	var header []string

	// find users first, and tweak description to add missing users
	reporterID := giteaFindUserID(reporter)
	if reporterID == -1 {
		header = append(header, fmt.Sprintf("    Originally reported by %s", reporter))
		reporterID = defaultAuthorID
	}
	var ownerID sql.NullString
	if owner != "" {
		tmp := giteaFindUserID(owner)
		if tmp == -1 {
			header = append(header, fmt.Sprintf("    Originally assigned to %s", owner))
			ownerID.String = fmt.Sprintf("%d", defaultAssigneeID)
			ownerID.Valid = true
		} else {
			ownerID.String = fmt.Sprintf("%d", tmp)
			ownerID.Valid = true
		}
	} else {
		ownerID.Valid = false
	}
	if len(header) > 0 {
		description = fmt.Sprintf("%s\n\n%s", strings.Join(header, "\n"), description)
	}

	issueID := giteaAddIssue(ticketID, summary, reporterID, milestone, ownerID, owner, closed, description, created)

	return issueID
}

func importTicketLabels(issueID int64, component string, severity string, priority string, version string, resolution string, typ string) {
	var lbl string
	if component != "" {
		lbl = "Component / " + component
		giteaAddIssueLabel(issueID, lbl)
	}

	if severity != "" {
		lbl = "Severity / " + severity
		giteaAddIssueLabel(issueID, lbl)
	}

	if priority != "" {
		lbl = "Priority / " + priority
		giteaAddIssueLabel(issueID, lbl)
	}

	if version != "" {
		lbl = "Version / " + version
		giteaAddIssueLabel(issueID, lbl)
	}

	if resolution != "" {
		lbl = "Resolution / " + resolution
		giteaAddIssueLabel(issueID, lbl)
	}

	if typ != "" {
		lbl = "Type / " + typ
		giteaAddIssueLabel(issueID, lbl)
	}
}

func importTicketComment(issueID int64, ticketID int64, time int64, author, comment string) int64 {
	var header []string

	prefix := fmt.Sprintf("ticket:%d", ticketID)
	comment = tracTextToMarkdown(comment, prefix)	// trac format

	// find users first, and tweak description to add missing users
	authorID := giteaFindUserID(author)
	if authorID == -1 {
		header = append(header, fmt.Sprintf("    Original comment by %s", author))
		authorID = defaultAuthorID
	}
	if len(header) > 0 {
		comment = fmt.Sprintf("%s\n\n%s", strings.Join(header, "\n"), comment)
	}

	return giteaAddComment(issueID, authorID, comment, time)
}

func importTicketAttachment(issueID int64, ticketID int64, time int64, size int64, author string, fname string, desc string) string {
	comment := fmt.Sprintf("**Attachment** %s (%d bytes) added\n\n%s",
		fname, size, desc)
	commentID := importTicketComment(issueID, ticketID, time, author, comment)

	tracPath := tracAttachmentPath(ticketID, fname)
	_, err := os.Stat(tracPath)
	if err != nil {
		log.Fatal(err)
	}
	elems := strings.Split(tracPath, "/")
	tracDir := elems[len(elems)-2]
	tracFile := elems[len(elems)-1]

	// 78ac is l33t for trac (horrible, I know)
	uuid := fmt.Sprintf("000078ac-%s-%s-%s-%s",
		tracDir[0:4], tracDir[4:8], tracDir[8:12],
		tracFile[0:12])

	// TODO: use a different uuid if file exists ?
	// TODO: avoid inserting record if uuid exist !
	giteaAddAttachment(uuid, issueID, commentID, fname, time)

	giteaPath := giteaAttachmentPath(uuid)
	giteaCopyFile(tracPath, giteaPath)

	return uuid
}

func importTicketAttachments(id int64, issueID int64, created int64) int64 {
	rows, err := tracDB.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, filename, description, size
			FROM attachment
  			WHERE type = 'ticket' AND id = $1
			ORDER BY time asc`, id)
	if err != nil {
		log.Fatal(err)
	}

	lastUpdate := created
	for rows.Next() {
		var time, size int64
		var author, fname, desc string
		if err := rows.Scan(&time, &author, &fname, &desc, &size); err != nil {
			log.Fatal(err)
		}

		fmt.Println(" adding attachment by", author)
		if lastUpdate > time {
			lastUpdate = time
		}
		importTicketAttachment(issueID, id, time, size, author, fname, desc)
	}

	return lastUpdate
}

func importTicketComments(ticketID int64, issueID int64, lastUpdate int64) {
	rows, err := tracDB.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, COALESCE(newvalue, '') newval
			FROM ticket_change where ticket = $1 AND field = 'comment' AND trim(COALESCE(newvalue, ''), ' ') != ''
			ORDER BY time asc`, ticketID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var time int64
		var author, comment string
		if err := rows.Scan(&time, &author, &comment); err != nil {
			log.Fatal(err)
		}
		fmt.Println(" adding comment by", author)
		if lastUpdate > time {
			lastUpdate = time
		}
		importTicketComment(issueID, ticketID, time, author, comment)
	}

	// Update issue modification time
	giteaSetIssueUpdateTime(issueID, lastUpdate)
}

func importTickets() {
	// NOTE: trac timestamps are to the microseconds, we just need seconds
	rows, err := tracDB.Query(`
		SELECT
			t.id,
			t.type,
			CAST(t.time*1e-6 AS int8),
			COALESCE(t.component, ''),
			COALESCE(t.severity,''),
			COALESCE(t.priority,''),
			COALESCE(t.owner,''),
			t.reporter,
			COALESCE(t.version,''),
			COALESCE(t.milestone,''),
			lower(COALESCE(t.status, '')),
			COALESCE(t.resolution,''),
			COALESCE(t.summary, ''),
			COALESCE(t.description, '')
		FROM ticket t ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	closedCount := 0
	for rows.Next() {
		var ticketID, created int64
		var component, ticketType, severity, priority, owner, reporter, version, milestone, status, resolution, summary, description string
		if err := rows.Scan(&ticketID, &ticketType, &created, &component, &severity, &priority, &owner, &reporter,
			&version, &milestone, &status, &resolution, &summary, &description); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Adding ticket", ticketID, " - ", summary)
		count++
		closed := status == "closed"
		if closed {
			closedCount++
		}

		issueID := importTicket(ticketID, created, owner, reporter, milestone, closed, summary, description)
		importTicketLabels(issueID, component, severity, priority, version, resolution, ticketType)
		lastUpdate := importTicketAttachments(ticketID, issueID, created)
		importTicketComments(ticketID, issueID, lastUpdate)
	}

	giteaUpdateRepoIssueCount(count, closedCount);

	// TODO: Update issue count for new labels
}

func main() {
	parseArgs();
	validateArgs();

	tracDB = getTracDB()
	giteaDB = getGiteaDB();

	giteaRepoID = giteaFindRepoID(giteaUserName, giteaRepoName)
	adminUser := giteaFindAdminUserName()
	adminUserID := giteaFindUserID(adminUser)
	defaultAssigneeID = giteaFindUserOrAdminID(defaultAssignee, adminUserID)
	defaultAuthorID = giteaFindUserOrAdminID(defaultAuthor, adminUserID)

	if !wikiOnly {
		importComponents();
		importPriorities();
		importSeverities();
		importVersions();
		importTypes();
		importResolutions();
		importMilestones();
		importTickets();
	}

	if ! dbOnly {
		importWiki();
	}
}
