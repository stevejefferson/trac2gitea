package markdown

import (
	"regexp"
	"strconv"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
)

// regexps for bracketted trac links
var doubleBracketImageLinkRegexp = regexp.MustCompile(`\[\[Image\(([^,\)]+)[^\]]*\]\]`)
var doubleBracketLinkRegexp = regexp.MustCompile(`\[\[([[:alpha:]][^|\]]*)(?:\|([^\]]+))?\]\]`)
var singleBracketLinkRegexp = regexp.MustCompile(`\[([[:alpha:]][^ \]]*)(?: +([^\]]+))?\]`) // no need to exclude '[' before initial '[' - Trac image and double bracket links get converted before this

// regexps for "marked" links - these are links which have been converted from their Trac form into an intermediate ("marked") form prior to final conversion to markdown
var noTextMarkedLinkRegexp = regexp.MustCompile(`((?:[^!]\[\])|[^\]])\(@@([^@]+)@@\)`)
var textMarkedLinkRegexp = regexp.MustCompile(`\(@@([^@]+)@@\)`)

// regexps for trac-supported links
var httpLinkRegexp = regexp.MustCompile(`https?://[^[:space:]]+`)
var htdocsLinkRegexp = regexp.MustCompile(`htdocs:([^[:space:]]+)`)
var wikiLinkRegexp = regexp.MustCompile(`(?:wiki:((?:[[:upper:]][[:lower:]]*)+))|((?:[[:upper:]][[:lower:]]+){2,})`) // Trac camel case with 'wiki:' prefix is more lax than without prefix
var ticketLinkRegexp = regexp.MustCompile(`ticket:([[:digit:]]+)`)
var ticketCommentLinkRegexp = regexp.MustCompile(`comment:([[:digit:]]+):ticket:([[:digit:]]+)`)
var milestoneLinkRegexp = regexp.MustCompile(`milestone:([^[:space:]]+)`)
var attachmentLinkRegexp = regexp.MustCompile(`attachment:([^[:space:]]+)[^:]`) // exclude trailing colon to avoid confusion with specific ticket attachments
var ticketAttachmentLinkRegexp = regexp.MustCompile(`attachment:([^[:space:]]+):ticket:([[:digit:]]+)`)
var changesetLinkRegexp = regexp.MustCompile(`changeset:"([[:alnum:]]+)/[^["]]+"`)
var sourceLinkRegexp = regexp.MustCompile(`source:"[^"/]+/([^"]+)"`)

func (converter *DefaultConverter) resolveHTTPLink(link string) string {
	return link // http links require no additional processing
}

func (converter *DefaultConverter) resolveHtdocsLink(link string) string {
	// any htdocs file needs copying from trac htdocs directory to an equivalent wiki subdirectory
	htdocsPath := htdocsLinkRegexp.ReplaceAllString(link, `$1`)
	tracHtdocsPath := converter.tracAccessor.GetFullPath("htdocs", htdocsPath)
	wikiHtdocsRelPath := "htdocs/" + htdocsPath
	wikiHtdocsURL := converter.wikiAccessor.CopyFile(tracHtdocsPath, wikiHtdocsRelPath)
	return wikiHtdocsURL
}

func (converter *DefaultConverter) resolveWikiLink(link string) string {
	wikiPageName := wikiLinkRegexp.ReplaceAllString(link, `$1`) // 'wiki:' prefix case
	if wikiPageName == "" {
		wikiPageName = wikiLinkRegexp.ReplaceAllString(link, `$2`) // unprefixed case
	}

	translatedPageName := converter.wikiAccessor.TranslatePageName(wikiPageName)
	return translatedPageName
}

func (converter *DefaultConverter) resolveTicketLink(link string) string {
	ticketIDStr := ticketLinkRegexp.ReplaceAllString(link, `$1`)
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		log.Warnf("found invalid Trac ticket reference %s\n" + link)
		return link
	}

	// validate ticket id
	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link
	}

	ticketLink := "#" + ticketIDStr
	return ticketLink
}

func (converter *DefaultConverter) resolveTicketCommentLink(link string) string {
	ticketCommentNumStr := ticketCommentLinkRegexp.ReplaceAllString(link, `$1`)
	var ticketCommentNum int64
	ticketCommentNum, err := strconv.ParseInt(ticketCommentNumStr, 10, 64)
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
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link
	}

	// find gitea ID for trac comment
	// - unfortunately the only real linkage between the trac comment number and gitea comment id here is the comment string itself
	commentStr := converter.tracAccessor.GetCommentString(ticketID, ticketCommentNum)
	commentID := converter.giteaAccessor.GetCommentID(issueID, commentStr)

	commentURL := converter.giteaAccessor.GetCommentURL(issueID, commentID)

	return commentURL
}

func (converter *DefaultConverter) resolveMilestoneLink(link string) string {
	milestoneName := milestoneLinkRegexp.ReplaceAllString(link, `$1`)
	milestoneID := converter.giteaAccessor.GetMilestoneID(milestoneName)
	if milestoneID == -1 {
		log.Warnf("cannot find milestone \"%s\" referenced by Trac link \"%s\"\n", milestoneName, link)
		return link
	}

	milestoneURL := converter.giteaAccessor.GetMilestoneURL(milestoneID)
	return milestoneURL
}

func (converter *DefaultConverter) resolveNamedAttachmentLink(link string, ticketID int64, attachmentName string) string {
	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link
	}

	uuid := converter.giteaAccessor.GetAttachmentUUID(issueID, attachmentName)
	if uuid == "" {
		log.Warnf("cannot find attachment \"%s\" for issue %d referenced by Trac link \"%s\"\n", attachmentName, issueID, link)
		return link
	}

	attachmentURL := converter.giteaAccessor.GetAttachmentURL(uuid)
	return attachmentURL
}

func (converter *DefaultConverter) resolveAttachmentLink(link string) string {
	attachmentName := attachmentLinkRegexp.ReplaceAllString(link, `$1`)
	return converter.resolveNamedAttachmentLink(link, converter.ticketID, attachmentName)
}

func (converter *DefaultConverter) resolveTicketAttachmentLink(link string) string {
	attachmentName := ticketAttachmentLinkRegexp.ReplaceAllString(link, `$1`)
	ticketIDStr := ticketAttachmentLinkRegexp.ReplaceAllString(link, `$2`)
	var ticketID int64
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return converter.resolveNamedAttachmentLink(link, ticketID, attachmentName)
}

func (converter *DefaultConverter) resolveChangesetLink(link string) string {
	changesetID := changesetLinkRegexp.ReplaceAllString(link, `$1`)
	changesetURL := converter.giteaAccessor.GetCommitURL(changesetID)
	return changesetURL
}

func (converter *DefaultConverter) resolveSourceLink(link string) string {
	sourcePath := sourceLinkRegexp.ReplaceAllString(link, `$1`)
	sourceURL := converter.giteaAccessor.GetSourceURL("master", sourcePath) // AFAICT Trac source URL does not include the git branch so we'll assume "master"
	return sourceURL
}

// convertBrackettedTracLinks converts the various forms of (square) bracketted Trac links into an unbracketted form.
// The conversion performed here is partial: this method is solely responsible for disposing of the Trac bracketting
// - any resolution of actual trac links is done later
func (converter *DefaultConverter) convertBrackettedTracLinks(in string) string {
	out := in

	out = doubleBracketImageLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// convert Trac image links to markdown but leave the link unprocessed
		// - it will get dealt with later
		link := doubleBracketImageLinkRegexp.ReplaceAllString(match, "$1")
		return "![]" + link
	})

	out = doubleBracketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// convert Trac double bracket links into Trac single bracket links
		// - if we convert directly to markdown here, the "[<text>]" part of the markdown will get misinterpreted as a Trac single bracket link
		link := doubleBracketLinkRegexp.ReplaceAllString(match, "$1")
		text := doubleBracketLinkRegexp.ReplaceAllString(match, "$2")
		if text == "" {
			return "[" + link + "]"
		}

		return "[" + link + " " + text + "]"
	})

	out = singleBracketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// convert Trac single bracket links to markdown but leave the link unprocessed
		// - it will get dealt with later
		link := singleBracketLinkRegexp.ReplaceAllString(match, "$1")
		text := singleBracketLinkRegexp.ReplaceAllString(match, "$2")
		return "[" + text + "]" + link
	})

	return out
}

// convertUnbrackettedTracLinks converts Trac-style links after any surrounding Trac bracketting and link texts have been processed
func (converter *DefaultConverter) convertUnbrackettedTracLinks(in string) string {
	out := in

	out = httpLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		httpLink := converter.resolveHTTPLink(match)
		return markLink(httpLink)
	})

	out = htdocsLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		htdocsLink := converter.resolveHtdocsLink(match)
		return markLink(htdocsLink)
	})

	out = wikiLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		wikiLink := converter.resolveWikiLink(match)
		return markLink(wikiLink)
	})

	out = ticketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		ticketLink := converter.resolveTicketLink(match)
		return markLink(ticketLink)
	})

	out = ticketCommentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		ticketCommentLink := converter.resolveTicketCommentLink(match)
		return markLink(ticketCommentLink)
	})

	out = milestoneLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		milestoneLink := converter.resolveMilestoneLink(match)
		return markLink(milestoneLink)
	})

	out = attachmentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		attachmentLink := converter.resolveAttachmentLink(match)
		return markLink(attachmentLink)
	})

	out = ticketAttachmentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		ticketAttachmentLink := converter.resolveTicketAttachmentLink(match)
		return markLink(ticketAttachmentLink)
	})

	out = changesetLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		changesetLink := converter.resolveChangesetLink(match)
		return markLink(changesetLink)
	})

	out = sourceLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		sourceLink := converter.resolveSourceLink(match)
		return markLink(sourceLink)
	})

	return out
}

// markLink marks the URL part of our links with a unique marker so that they can be picked up by later converters.
// We cannot just convert to markdown at this stage because markdown's round brackets are insufficiently unique
// and would result in the later converters misidentifying normal bracketted text as links.
// Hence we put a marker in here and later convert that marker to the final markdown.
func markLink(in string) string {
	return "(@@" + in + "@@)"
}

// unmarkLinks removes the "marking" placed around links by markLinks and converts them into their final markdown format
func (converter *DefaultConverter) unmarkLinks(in string) string {
	out := in
	out = noTextMarkedLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// (marked) links with no accompanying comment are converted into markdown "automatic" links
		leadingChars := noTextMarkedLinkRegexp.ReplaceAllString(match, `$1`)
		markdownURL := noTextMarkedLinkRegexp.ReplaceAllString(match, `$2`)

		// need to replace any leading chars except for any square brackets forming the empty link text
		prefix := strings.Replace(leadingChars, "[]", "", 1)
		return prefix + "<" + markdownURL + ">"
	})

	out = textMarkedLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// any remaining (marked) links must have an accompanying comment so are converted into normal markdown links
		markdownURL := textMarkedLinkRegexp.ReplaceAllString(match, `$1`)
		return "(" + markdownURL + ")"
	})

	return out
}

func (converter *DefaultConverter) convertLinks(in string) string {
	out := in

	// conversion occurs in three distinct phases with each phase dealing with one part of the link syntax
	// and leaving the remainder for the next stage
	out = converter.convertBrackettedTracLinks(out)
	out = converter.convertUnbrackettedTracLinks(out)
	out = converter.unmarkLinks(out)
	return out
}

var httpLinkDisguiseRegexp = regexp.MustCompile(`(https?)://`)
var httpLinkUndisguiseRegexp = regexp.MustCompile(`(https?):@@`)

// disguiseLinks temporarily disguises links into a format that doesn't interfere with other Trac -> markdown regexps
// - in particular the '//' in 'http(s)://...' clashes with Trac's '//' italics marker
func (converter *DefaultConverter) disguiseLinks(in string) string {
	out := in
	out = httpLinkDisguiseRegexp.ReplaceAllString(out, `$1:@@`)
	return out
}

// undisguiseLink converts temporarily "disguised" links back to their correct format.
func (converter *DefaultConverter) undisguiseLinks(in string) string {
	out := in
	out = httpLinkUndisguiseRegexp.ReplaceAllString(out, `$1://`)
	return out
}
