package markdown

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
)

func markdownAutomaticLink(link string) string {
	return "<" + link + ">"
}

func markdownImageLink(link string) string {
	return "![](" + link + ")"
}

func markdownLink(link string, text string) string {
	return "[" + text + "](" + link + ")"
}

// regexps for bracketted trac links
var doubleBracketImageLinkRegexp = regexp.MustCompile(`\[\[Image\(([^,\)]+)[^\]]*\]\]`)
var doubleBracketLinkRegexp = regexp.MustCompile(`\[\[([[:alpha:]][^|]*)\|([^\]]+\])\]`)
var singleBracketLinkRegexp = regexp.MustCompile(`([^\[])\[([[:alpha:]][^ \]]*) +([^\]\n]+)\]`) // exclude '[' before initial '[' to avoid picking up '[[Image...'

// regexp for text which might be an unbracketted trac link
var potentialUnbrackettedLinkRegexp = regexp.MustCompile(`[[:space:]]([[:alpha:]][^[:space:]\]]*)[[:space:]]`)

// regexps for trac-supported links
var httpLinkRegexp = regexp.MustCompile(`^https?://[^[:space:]]+`)
var htdocsLinkRegexp = regexp.MustCompile(`^htdocs:([^[:space:]]+)`)
var wikiCamelCaseLinkRegexp = regexp.MustCompile(`^((?:[[:upper:]][[:lower:]]+){2,})`)
var wikiLinkRegexp = regexp.MustCompile(`^wiki:((?:[[:upper:]][[:lower:]]*)+)`) // 'wiki:' prefix allows camel case to be more lax than above

var ticketLinkRegexp = regexp.MustCompile(`^ticket:([[:digit:]]+)`)
var ticketCommentLinkRegexp = regexp.MustCompile(`^comment:([[:digit:]]+):ticket:([[:digit:]]+)`)
var milestoneLinkRegexp = regexp.MustCompile(`^milestone:([^[:space:]]+)`)
var attachmentLinkRegexp = regexp.MustCompile(`^attachment:([^[:space:]]+)[^:]`) // exclude trailing colon to avoid confusion with specific ticket attachments
var ticketAttachmentLinkRegexp = regexp.MustCompile(`^attachment:([^[:space:]]+):ticket:([[:digit:]]+)`)
var changesetLinkRegexp = regexp.MustCompile(`^changeset:"([[:alnum:]]+)/[^["]]+"`)
var sourceLinkRegexp = regexp.MustCompile(`^source:"[^"/]+/([^"]+)"`)

func (converter *Converter) resolveHTTPLink(link string) string {
	return link // http* links are returned as is
}

func (converter *Converter) resolveHtdocsLink(link string) string {
	// any htdocs file needs copying from trac htdocs directory to an equivalent wiki subdirectory
	htdocsPath := htdocsLinkRegexp.ReplaceAllString(link, `$1`)
	tracHtdocsPath := filepath.Join(converter.tracAccessor.RootDir, "htdocs", htdocsPath)
	wikiHtdocsRelPath := "htdocs/" + htdocsPath
	converter.wikiAccessor.CopyFile(tracHtdocsPath, wikiHtdocsRelPath)
	return "../raw/" + wikiHtdocsRelPath // htodcs subdirectory should be referenceable via Gitea "raw" repo path...
}

func (converter *Converter) resolveWikiCamelCaseLink(link string) string {
	wikiPageName := wikiCamelCaseLinkRegexp.ReplaceAllString(link, `$1`)
	translatedPageName := converter.wikiAccessor.TranslatePageName(wikiPageName)
	return translatedPageName
}

func (converter *Converter) resolveWikiLink(link string) string {
	wikiPageName := wikiLinkRegexp.ReplaceAllString(link, `$1`)
	translatedPageName := converter.wikiAccessor.TranslatePageName(wikiPageName)
	return translatedPageName
}

func (converter *Converter) resolveTicketLink(link string) string {
	return ticketLinkRegexp.ReplaceAllString(link, `#$1`) // convert into '#nnn' ticket link
}

func (converter *Converter) resolveTicketCommentLink(link string) string {
	commentIDStr := ticketCommentLinkRegexp.ReplaceAllString(link, `$1`)
	var commentID int64
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	ticketIDStr := ticketCommentLinkRegexp.ReplaceAllString(link, `$2`)
	var ticketID int64
	ticketID, err = strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Printf("Warning: cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link
	}
	commentURL := converter.giteaAccessor.GetCommentURL(issueID, commentID)

	return commentURL
}

func (converter *Converter) resolveMilestoneLink(link string) string {
	milestoneName := milestoneLinkRegexp.ReplaceAllString(link, `$1`)
	milestoneID := converter.giteaAccessor.GetMilestoneID(milestoneName)
	if milestoneID == -1 {
		log.Printf("Warning: cannot find milestone \"%s\" referenced by Trac link \"%s\"\n", milestoneName, link)
		return link
	}

	milestoneURL := converter.giteaAccessor.GetMilestoneURL(milestoneID)
	return milestoneURL
}

func (converter *Converter) resolveNamedAttachmentLink(link string, ticketID int64, attachmentName string) string {
	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Printf("Warning: cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link
	}

	uuid := converter.giteaAccessor.GetAttachmentUUID(issueID, attachmentName)
	if uuid == "" {
		log.Printf("Warning: cannot find attachment \"%s\" for issue %d referenced by Trac link \"%s\"\n", attachmentName, issueID, link)
		return link
	}

	return converter.giteaAccessor.GetAttachmentURL(uuid)
}

func (converter *Converter) resolveAttachmentLink(link string) string {
	attachmentName := attachmentLinkRegexp.ReplaceAllString(link, `$1`)
	return converter.resolveNamedAttachmentLink(link, converter.ticketID, attachmentName)
}

func (converter *Converter) resolveTicketAttachmentLink(link string) string {
	attachmentName := ticketAttachmentLinkRegexp.ReplaceAllString(link, `$1`)
	ticketIDStr := ticketAttachmentLinkRegexp.ReplaceAllString(link, `$2`)
	var ticketID int64
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return converter.resolveNamedAttachmentLink(link, ticketID, attachmentName)
}

func (converter *Converter) resolveChangesetLink(link string) string {
	changesetID := changesetLinkRegexp.ReplaceAllString(link, `$1`)
	return converter.giteaAccessor.GetCommitURL(changesetID)
}

func (converter *Converter) resolveSourceLink(link string) string {
	sourcePath := sourceLinkRegexp.ReplaceAllString(link, `$1`)
	return converter.giteaAccessor.GetSourceURL("master", sourcePath) // AFAICT Trac source URL does not include the git branch so we'll assume "master"
}

// resolveLink resolves a Trac-style link into a Gitea link suitable for embedding in Markdown
// returns "" if provided link is not recognised as a Trac link
func (converter *Converter) resolveLink(in string) string {
	if httpLinkRegexp.MatchString(in) {
		return converter.resolveHTTPLink(in)
	}

	if htdocsLinkRegexp.MatchString(in) {
		return converter.resolveHtdocsLink(in)
	}

	if wikiLinkRegexp.MatchString(in) {
		return converter.resolveWikiLink(in)
	}

	if wikiCamelCaseLinkRegexp.MatchString(in) {
		return converter.resolveWikiCamelCaseLink(in)
	}

	if ticketLinkRegexp.MatchString(in) {
		return converter.resolveTicketLink(in)
	}

	if ticketCommentLinkRegexp.MatchString(in) {
		return converter.resolveTicketCommentLink(in)
	}

	if milestoneLinkRegexp.MatchString(in) {
		return converter.resolveMilestoneLink(in)
	}

	if attachmentLinkRegexp.MatchString(in) {
		return converter.resolveAttachmentLink(in)
	}

	if ticketAttachmentLinkRegexp.MatchString(in) {
		return converter.resolveTicketAttachmentLink(in)
	}

	if changesetLinkRegexp.MatchString(in) {
		return converter.resolveChangesetLink(in)
	}

	if sourceLinkRegexp.MatchString(in) {
		return converter.resolveSourceLink(in)
	}

	return ""
}

func (converter *Converter) convertLinks(in string) string {
	out := in

	out = doubleBracketImageLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		link := doubleBracketImageLinkRegexp.ReplaceAllString(match, "$1")
		resolvedLink := converter.resolveLink(link)
		if resolvedLink == "" {
			return match
		}
		return markdownImageLink(resolvedLink)
	})

	out = doubleBracketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		link := doubleBracketLinkRegexp.ReplaceAllString(match, "$1")
		text := doubleBracketLinkRegexp.ReplaceAllString(match, "$2")
		resolvedLink := converter.resolveLink(link)
		if resolvedLink == "" {
			return match
		}
		return markdownLink(resolvedLink, text)
	})

	out = singleBracketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		leadingChar := singleBracketLinkRegexp.ReplaceAllString(match, "$1")
		link := singleBracketLinkRegexp.ReplaceAllString(match, "$2")
		text := singleBracketLinkRegexp.ReplaceAllString(match, "$3")
		resolvedLink := converter.resolveLink(link)
		if resolvedLink == "" {
			return match
		}
		return leadingChar + markdownLink(resolvedLink, text)
	})

	out = potentialUnbrackettedLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		link := potentialUnbrackettedLinkRegexp.ReplaceAllString(match, "$1")
		resolvedLink := converter.resolveLink(link)
		if resolvedLink == "" {
			return match
		}
		return markdownAutomaticLink(resolvedLink)
	})

	return out
}

var httpLinkDisguiseRegexp = regexp.MustCompile(`(https?)://`)
var httpLinkUndisguiseRegexp = regexp.MustCompile(`(https?):\|\|`)

// disguiseLinks temporarily disguises links into a format that doesn't interfere with other Trac -> markdown regexps
// - in particular the '//' in 'http(s)://...' clashes with Trac's '//' italics marker
func (converter *Converter) disguiseLinks(in string) string {
	out := in
	out = httpLinkDisguiseRegexp.ReplaceAllString(out, `$1:||`)
	return out
}

// undisguiseLink converts temporarily "disguised" links back to their correct format.
func (converter *Converter) undisguiseLinks(in string) string {
	out := in
	out = httpLinkUndisguiseRegexp.ReplaceAllString(out, `$1://`)
	return out
}
