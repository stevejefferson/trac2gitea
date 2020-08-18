package markdown

import (
	"regexp"
	"strconv"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
)

var (
	// regexp for trac '[[Image(<image>...,link=<link>)]]': $1=image, $2=link
	doubleBracketImageLinkRegexp = regexp.MustCompile(`\[\[Image\(([^,)]+)(?:, *link=([[:alnum:]\-._~:/?#@!$&'(*+;%=]+))?[^\]]*\]\]`)

	// regexp for trac '[[<link>]]' and '[[<link>|<text>]]': $1=link, $2=text
	doubleBracketLinkRegexp = regexp.MustCompile(`\[\[([[:alpha:]][^|\]]*)(?:\|([^\]]+))?\]\]`)

	// regexp for trac '[<link>]' and '[<link> <text>]': $1=link, $2=text
	// note: trac image and double bracket links are processed before this so we do not need to exclude a leading '[' in the regexp
	singleBracketLinkRegexp = regexp.MustCompile(`\[([[:alpha:]][^ \]]*)(?: +([^\]]+))?\]`)

	// regexp for 'http://...' and 'https://...' links
	httpLinkRegexp = regexp.MustCompile(`https?://[[:alnum:]\-._~:/?#@!$&'()*+,;%=]+`)

	// regexp for trac 'htdocs:<link>': $1=link
	htdocsLinkRegexp = regexp.MustCompile(`htdocs:([[:alnum:]\-._~:/?#@!$&'()*+,;%=]+)`)

	// regexp for trac '<CamelCase>' wiki links: $1=leading char, $2=CamelCase
	// note: leading char (if any) must be a space or ']'
	//       - a space constitutes a "start of word" for an unbracketted CamelCase link,
	//       - a ']' constitutes the end of the link comment after conversion of the various trac bracketting syntaxes above
	wikiCamelCaseLinkRegexp = regexp.MustCompile(`([[:space:]\]]|\A)((?:[[:upper:]][[:lower:]]+){2,})`)

	// regexp for trac 'wiki:<CamelCase>' links: $1=leading text, $2=CamelCase
	// notes: 1. rules on what constitutes "CamelCase" are more lax than for plain <CamelCase> variant
	//        2. leading char is used to exclude a leading ':' and so avoid confusion with wiki page attachment links
	wikiLinkRegexp = regexp.MustCompile(`([^:])wiki:((?:[[:upper:]][[:lower:]]*)+)`)

	// regexp for a trac 'ticket:<ticketID>' link: $1=leading char, $2=ticketID
	// note: leading char is used to exclude a leading ':' and so avoid confusion with ticket comment and attachment links
	ticketLinkRegexp = regexp.MustCompile(`([^:])ticket:([[:digit:]]+)`)

	// regexp for a trac  'comment:<commentNum>:ticket:<ticketID>' link: $1=commentNum, $2=ticketID
	ticketCommentLinkRegexp = regexp.MustCompile(`comment:([[:digit:]]+):ticket:([[:digit:]]+)`)

	// regexp for a trac 'milestone:<milestoneName>' link: $1=milestoneName
	milestoneLinkRegexp = regexp.MustCompile(`milestone:([[:alnum:]\-._~:/?#@!$&'()*+,;%=]+)`)

	// regexp for a trac 'attachment:<attachmentName>' link: $1=attachmentName, $2=trailing char
	// note: trailing char is used to exclude a trailing ':' and so avoid confusion with explict wiki page and with ticket attachments
	attachmentLinkRegexp = regexp.MustCompile(`attachment:([[:alnum:]\-._~/?#@!$&'()*+,;%=]+)([^\:])`)

	// regexp for a trac 'attachment:<attachmentName>:wiki:<pageName>' link: $1=attachmentName, $2=pageName
	wikiAttachmentLinkRegexp = regexp.MustCompile(`attachment:([[:alnum:]\-._~/?#@!$&'()*+,;%=]+):wiki:((?:[[:upper:]][[:lower:]]*)+)`)

	// regexp for a trac 'attachment:<attachmentName>:ticket:<ticketID>' link: $1=attachmentName, $2=ticketID
	ticketAttachmentLinkRegexp = regexp.MustCompile(`attachment:([[:alnum:]\-._~/?#@!$&'()*+,;%=]+):ticket:([[:digit:]]+)`)

	// regexp for a trac 'changeset:<changesetID>' link: $1=commitID
	changesetLinkRegexp = regexp.MustCompile(`changeset:"([[:xdigit:]]+)/[^"]+"`)

	// regexp for a trac 'source:<sourcePath>' link: $1=sourcePath
	sourceLinkRegexp = regexp.MustCompile(`source:"[^/]+/([^"]+)"`)

	// regexp for recognising a "marked" link with no accompanying text: $1=leading chars, $2=link
	noTextMarkedLinkRegexp = regexp.MustCompile(`((?:[^!]\[\])|[^\]])\(@@([^@]+)@@\)`)

	// regexp for recognising a "marked" link with accompanying text: $1=link
	textMarkedLinkRegexp = regexp.MustCompile(`\(@@([^@]+)@@\)`)
)

// link resolution functions:
// 	These are responsible for extracting link information from its appropriate Trac link regexp and preparing that link for conversion to markdown.
//	The portion of the returned text corresponding to the link itself (as opposed to any extraneous characters that may have been hoovered up by the regexp)
//  should be "marked" using the markLink() function to identify it for later processing.
func (converter *DefaultConverter) resolveHTTPLink(link string) string {
	return markLink(link)
}

func (converter *DefaultConverter) resolveHtdocsLink(link string) string {
	// any htdocs file needs copying from trac htdocs directory to an equivalent wiki subdirectory
	htdocPath := htdocsLinkRegexp.ReplaceAllString(link, `$1`)
	tracHtdocPath := converter.tracAccessor.GetFullPath("htdocs", htdocPath)
	wikiHtdocRelPath := converter.giteaAccessor.GetWikiHtdocRelPath(htdocPath)
	converter.giteaAccessor.CopyFileToWiki(tracHtdocPath, wikiHtdocRelPath)
	wikiHtdocURL := converter.giteaAccessor.GetWikiFileURL(wikiHtdocRelPath)
	return markLink(wikiHtdocURL)
}

func (converter *DefaultConverter) resolveWikiCamelCaseLink(link string) string {
	leadingChar := wikiCamelCaseLinkRegexp.ReplaceAllString(link, `$1`)
	wikiPageName := wikiCamelCaseLinkRegexp.ReplaceAllString(link, `$2`)
	translatedPageName := converter.giteaAccessor.TranslateWikiPageName(wikiPageName)
	return leadingChar + markLink(translatedPageName)
}

func (converter *DefaultConverter) resolveWikiLink(link string) string {
	leadingChar := wikiLinkRegexp.ReplaceAllString(link, `$1`)
	wikiPageName := wikiLinkRegexp.ReplaceAllString(link, `$2`)
	translatedPageName := converter.giteaAccessor.TranslateWikiPageName(wikiPageName)
	return leadingChar + markLink(translatedPageName)
}

func (converter *DefaultConverter) resolveTicketLink(link string) string {
	leadingChar := ticketLinkRegexp.ReplaceAllString(link, `$1`)
	ticketIDStr := ticketLinkRegexp.ReplaceAllString(link, `$2`)
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		log.Warnf("found invalid Trac ticket reference %s\n" + link)
		return link // not a recognised link - do not mark
	}

	// validate ticket id
	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link // not a recognised link - do not mark
	}

	issueURL := converter.giteaAccessor.GetIssueURL(issueID)
	return leadingChar + markLink(issueURL)
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
		return link // not a recognised link - do not mark
	}

	// find gitea ID for trac comment
	// - unfortunately the only real linkage between the trac comment number and gitea comment id here is the comment string itself
	commentStr := converter.tracAccessor.GetCommentString(ticketID, ticketCommentNum)
	commentID := converter.giteaAccessor.GetCommentID(issueID, commentStr)

	commentURL := converter.giteaAccessor.GetCommentURL(issueID, commentID)
	return markLink(commentURL)
}

func (converter *DefaultConverter) resolveMilestoneLink(link string) string {
	milestoneName := milestoneLinkRegexp.ReplaceAllString(link, `$1`)
	milestoneID := converter.giteaAccessor.GetMilestoneID(milestoneName)
	if milestoneID == -1 {
		log.Warnf("cannot find milestone \"%s\" referenced by Trac link \"%s\"\n", milestoneName, link)
		return link // not a recognised link - do not mark
	}

	milestoneURL := converter.giteaAccessor.GetMilestoneURL(milestoneID)
	return markLink(milestoneURL)
}

func (converter *DefaultConverter) resolveAttachmentLink(link string) string {
	attachmentName := attachmentLinkRegexp.ReplaceAllString(link, `$1`)
	trailingChars := attachmentLinkRegexp.ReplaceAllString(link, `$2`)

	attachmentWikiRelPath := converter.giteaAccessor.GetWikiAttachmentRelPath(converter.wikiPage, attachmentName)
	attachmentWikiURL := converter.giteaAccessor.GetWikiFileURL(attachmentWikiRelPath)

	return markLink(attachmentWikiURL) + trailingChars
}

func (converter *DefaultConverter) resolveWikiAttachmentLink(link string) string {
	attachmentName := wikiAttachmentLinkRegexp.ReplaceAllString(link, `$1`)
	wikiPageName := wikiAttachmentLinkRegexp.ReplaceAllString(link, `$2`)

	attachmentWikiRelPath := converter.giteaAccessor.GetWikiAttachmentRelPath(wikiPageName, attachmentName)
	attachmentWikiURL := converter.giteaAccessor.GetWikiFileURL(attachmentWikiRelPath)

	return markLink(attachmentWikiURL)
}

func (converter *DefaultConverter) resolveTicketAttachmentLink(link string) string {
	attachmentName := ticketAttachmentLinkRegexp.ReplaceAllString(link, `$1`)
	ticketIDStr := ticketAttachmentLinkRegexp.ReplaceAllString(link, `$2`)
	var ticketID int64
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	issueID := converter.giteaAccessor.GetIssueID(ticketID)
	if issueID == -1 {
		log.Warnf("cannot find Gitea issue for ticket %d referenced by Trac link \"%s\"\n", ticketID, link)
		return link // not a recognised link - do not mark
	}

	uuid := converter.giteaAccessor.GetAttachmentUUID(issueID, attachmentName)
	if uuid == "" {
		log.Warnf("cannot find attachment \"%s\" for issue %d referenced by Trac link \"%s\"\n", attachmentName, issueID, link)
		return link // not a recognised link - do not mark
	}

	attachmentURL := converter.giteaAccessor.GetAttachmentURL(uuid)
	return markLink(attachmentURL)
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
		image := doubleBracketImageLinkRegexp.ReplaceAllString(match, "$1")
		link := doubleBracketImageLinkRegexp.ReplaceAllString(match, "$2")
		if link == "" {
			return "![]" + image
		}
		return "[![]" + image + "]" + link
	})

	out = doubleBracketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// convert Trac double bracket links into Trac single bracket links
		// - if we convert directly to markdown here, the "[<text>]" part of the markdown will get misinterpreted as a Trac single bracket link
		link := doubleBracketLinkRegexp.ReplaceAllString(match, "$1")
		text := doubleBracketLinkRegexp.ReplaceAllString(match, "$2")

		if text == "" {
			// special case: '[[br]]' is a page break in Trac and is dealt with elsewhere
			if strings.EqualFold(link, "br") {
				return match
			}
			return "[" + link + "]"
		}

		return "[" + link + " " + text + "]"
	})

	out = singleBracketLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// convert Trac single bracket links to markdown but leave the link unprocessed; it will get dealt with later
		link := singleBracketLinkRegexp.ReplaceAllString(match, "$1")
		text := singleBracketLinkRegexp.ReplaceAllString(match, "$2")

		// special case: '[br]' can be assumed to be the inner section of a '[[br]]' and is a page break in Trac which is dealt with elsewhere
		if text == "" && strings.EqualFold(link, "br") {
			return match
		}

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

	out = wikiAttachmentLinkRegexp.ReplaceAllStringFunc(out, func(match string) string {
		return converter.resolveWikiAttachmentLink(match)
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
