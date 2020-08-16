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
var wikiCamelCaseLinkRegexp = regexp.MustCompile(`([^:]|\A)((?:[[:upper:]][[:lower:]]+){2,})`) // exclude leading colon to avoid confusion with 'wiki:...'
var wikiLinkRegexp = regexp.MustCompile(`wiki:((?:[[:upper:]][[:lower:]]*)+)`)                 // 'wiki:' prefix allows camel case to be more lax than above
var ticketLinkRegexp = regexp.MustCompile(`ticket:([[:digit:]]+)`)
var ticketCommentLinkRegexp = regexp.MustCompile(`comment:([[:digit:]]+):ticket:([[:digit:]]+)`)
var milestoneLinkRegexp = regexp.MustCompile(`milestone:([^[:space:]]+)`)
var attachmentLinkRegexp = regexp.MustCompile(`attachment:([^[:space:]]+)[^:]`) // exclude trailing colon to avoid confusion with specific ticket attachments
var ticketAttachmentLinkRegexp = regexp.MustCompile(`attachment:([^[:space:]]+):ticket:([[:digit:]]+)`)
var changesetLinkRegexp = regexp.MustCompile(`changeset:"([[:alnum:]]+)/[^["]]+"`)
var sourceLinkRegexp = regexp.MustCompile(`source:"[^"/]+/([^"]+)"`)

// resolution for various Trac links
// The links returned by these functions should be "marked" using 'markLink' to identify them during later processing

func (converter *DefaultConverter) resolveHTTPLink(link string) string {
	return markLink(link)
}

func (converter *DefaultConverter) resolveHtdocsLink(link string) string {
	// any htdocs file needs copying from trac htdocs directory to an equivalent wiki subdirectory
	htdocsPath := htdocsLinkRegexp.ReplaceAllString(link, `$1`)
	tracHtdocsPath := converter.tracAccessor.GetFullPath("htdocs", htdocsPath)
	wikiHtdocsRelPath := "htdocs/" + htdocsPath
	wikiHtdocsURL := converter.wikiAccessor.CopyFile(tracHtdocsPath, wikiHtdocsRelPath)
	return markLink(wikiHtdocsURL)
}

func (converter *DefaultConverter) resolveWikiCamelCaseLink(link string) string {
	leadingChar := wikiCamelCaseLinkRegexp.ReplaceAllString(link, `$1`)
	wikiPageName := wikiCamelCaseLinkRegexp.ReplaceAllString(link, `$2`)
	translatedPageName := converter.wikiAccessor.TranslatePageName(wikiPageName)
	return leadingChar + markLink(translatedPageName)
}

func (converter *DefaultConverter) resolveWikiLink(link string) string {
	wikiPageName := wikiLinkRegexp.ReplaceAllString(link, `$1`)
	translatedPageName := converter.wikiAccessor.TranslatePageName(wikiPageName)
	return markLink(translatedPageName)
}

func (converter *DefaultConverter) resolveTicketLink(link string) string {
	ticketLink := ticketLinkRegexp.ReplaceAllString(link, `#$1`) // convert into '#nnn' ticket link
	return markLink(ticketLink)
}

func (converter *DefaultConverter) resolveTicketCommentLink(link string) string {
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
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link // result has not been successfully identified as a link - do not mark
	}
	commentURL := converter.giteaAccessor.GetCommentURL(issueID, commentID)

	return markLink(commentURL)
}

func (converter *DefaultConverter) resolveMilestoneLink(link string) string {
	milestoneName := milestoneLinkRegexp.ReplaceAllString(link, `$1`)
	milestoneID := converter.giteaAccessor.GetMilestoneID(milestoneName)
	if milestoneID == -1 {
		log.Warnf("cannot find milestone \"%s\" referenced by Trac link \"%s\"\n", milestoneName, link)
		return link // result has not been successfully identified as a link - do not mark
	}

	milestoneURL := converter.giteaAccessor.GetMilestoneURL(milestoneID)
	return markLink(milestoneURL)
}

func (converter *DefaultConverter) resolveNamedAttachmentLink(link string, ticketID int64, attachmentName string) string {
	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link // result has not been successfully identified as a link - do not mark
	}

	uuid := converter.giteaAccessor.GetAttachmentUUID(issueID, attachmentName)
	if uuid == "" {
		log.Warnf("cannot find attachment \"%s\" for issue %d referenced by Trac link \"%s\"\n", attachmentName, issueID, link)
		return link // result has not been successfully identified as a link - do not mark
	}

	attachmentURL := converter.giteaAccessor.GetAttachmentURL(uuid)
	return markLink(attachmentURL)
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
	return markLink(changesetURL)
}

func (converter *DefaultConverter) resolveSourceLink(link string) string {
	sourcePath := sourceLinkRegexp.ReplaceAllString(link, `$1`)
	sourceURL := converter.giteaAccessor.GetSourceURL("master", sourcePath) // AFAICT Trac source URL does not include the git branch so we'll assume "master"
	return markLink(sourceURL)
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
		return converter.resolveHTTPLink(match)
	})

	out = htdocsLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveHtdocsLink(match)
	})

	// wikiCamelCaseLink must precede wikiLink otherwise wiki link will be stripped of its "wiki:" prefix then re-processed as camel case link
	// - of itself this is harmless but it can trip up mock expectations in the unit tests
	out = wikiCamelCaseLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveWikiCamelCaseLink(match)
	})

	out = wikiLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveWikiLink(match)
	})

	out = ticketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveTicketLink(match)
	})

	out = ticketCommentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveTicketCommentLink(match)
	})

	out = milestoneLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveMilestoneLink(match)
	})

	out = attachmentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveAttachmentLink(match)
	})

	out = ticketAttachmentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveTicketAttachmentLink(match)
	})

	out = changesetLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveChangesetLink(match)
	})

	out = sourceLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveSourceLink(match)
	})

	return out
}

// markLink marks the URL part of our links with a unique marker so that they can be picked up by later converters.
// We cannot just convert to markdown at this stage because markdown's round brackets are insufficiently unique
// and could result in the later converters misidentifying normal bracketted text as links.
// Hence we put a marker in here then convert that marker to the final markdown later on.
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
