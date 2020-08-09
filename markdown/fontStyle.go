package markdown

import "regexp"

// regexps for trac-supported font styles
var singleQuoteBoldRegexp = regexp.MustCompile(`'''((?:[^'\n][^'\n][^\n'])*[^!])'''`)
var singleQuoteItalicRegexp = regexp.MustCompile(`''((?:[^'\n][^'\n])*)''`)
var singleQuoteBoldItalicRegexp = regexp.MustCompile(`'''''((?:[^'\n][^'\n][^'\n][^'\n][^'\n])*)'''''`)
var doubleAsteriskBoldRegexp = regexp.MustCompile(`\*\*((?:[^*\n][^*\n])*)\*\*`)
var doubleSlashItalicRegexp = regexp.MustCompile(`//((?:[^/\n][^/\n])*)//`)
var underlineRegexp = regexp.MustCompile(`__((?:[^_\n][^_\n])*)__`)

// regexps for trac font styles for which we have no mappings
//var wikiCreoleStyleRegexp = regexp.MustCompile(`\*\*//!(.*)//\*\*`)
//var strikethroughRegexp = regexp.MustCompile(`~~(.*)~~`)
//var superscriptRegexp = regexp.MustCompile(`\^(.*)\^`)
//var subscriptRegexp = regexp.MustCompile(`,,(.*),,`)

// markdown replacement strings
var emphasisReplacementStr = `*$1*`
var strongReplacementStr = `**$1**`

func (converter *Converter) convertFontStyle(in string) string {
	out := in
	out = singleQuoteBoldItalicRegexp.ReplaceAllString(out, strongReplacementStr)
	out = singleQuoteBoldRegexp.ReplaceAllString(out, strongReplacementStr)
	out = singleQuoteItalicRegexp.ReplaceAllString(out, emphasisReplacementStr)
	out = doubleAsteriskBoldRegexp.ReplaceAllString(out, strongReplacementStr)
	out = doubleSlashItalicRegexp.ReplaceAllString(out, emphasisReplacementStr)
	out = underlineRegexp.ReplaceAllString(out, emphasisReplacementStr)

	return out
}
