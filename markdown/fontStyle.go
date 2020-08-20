// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

import (
	"regexp"
)

// regexps for trac-supported font styles
var singleQuoteBoldRegexp = regexp.MustCompile(`'''([^\n]+?)'''`)
var singleQuoteItalicRegexp = regexp.MustCompile(`''([^\n]+?)''`)
var singleQuoteBoldItalicRegexp = regexp.MustCompile(`'''''([^\n]+?)'''''`)
var doubleAsteriskBoldRegexp = regexp.MustCompile(`\*\*([^\n]+?)\*\*`)
var doubleSlashItalicRegexp = regexp.MustCompile(`//([^\n]+?)//`)
var underlineRegexp = regexp.MustCompile(`__([^\n]+?)__`)

// regexps for trac font styles for which we have no mappings
//var wikiCreoleStyleRegexp = regexp.MustCompile(`\*\*//!([^\n]+?)//\*\*`)
//var strikethroughRegexp = regexp.MustCompile(`~~([^\n]+?)~~`)
//var superscriptRegexp = regexp.MustCompile(`\^([^\n]+?)\^`)
//var subscriptRegexp = regexp.MustCompile(`,,([^\n]+?),,`)

// markdown replacement strings
var emphasisReplacementStr = `*$1*`
var strongReplacementStr = `**$1**`

func (converter *DefaultConverter) convertFontStyles(in string) string {
	out := in
	out = singleQuoteBoldItalicRegexp.ReplaceAllString(out, strongReplacementStr)
	out = singleQuoteBoldRegexp.ReplaceAllString(out, strongReplacementStr)
	out = singleQuoteItalicRegexp.ReplaceAllString(out, emphasisReplacementStr)
	out = doubleAsteriskBoldRegexp.ReplaceAllString(out, strongReplacementStr)
	out = doubleSlashItalicRegexp.ReplaceAllString(out, emphasisReplacementStr)
	out = underlineRegexp.ReplaceAllString(out, emphasisReplacementStr)

	return out
}
