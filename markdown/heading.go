// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

import (
	"regexp"
	"strings"
)

const maxHeadingLevel = 6

var headingRegexps [maxHeadingLevel]*regexp.Regexp
var headingAnchorRegexp = regexp.MustCompile(`(?m)^ +#([^[:space:]]+)`)

func compileRegexps() {
	tracHeadingDelimiter := ""
	for headingLevel := 0; headingLevel < maxHeadingLevel; headingLevel++ {
		tracHeadingDelimiter = tracHeadingDelimiter + "="

		// construct regexp for trac heading - this captures:
		// $1 = heading level delimiter (which we transform into the equivalent markdown delimiter)
		// $2 = heading text
		// $3 = any trailing text on line which may contain an optional Trac heading anchor
		// note: the trailing sequence of '='s  on trac headings turns out to be optional
		headingRegexpStr := `(?m)^(` + tracHeadingDelimiter + `) *([^=\n]+)(?:` + tracHeadingDelimiter + `)?(.*)$`
		headingRegexps[headingLevel] = regexp.MustCompile(headingRegexpStr)
	}
}

func init() {
	// pre-compile array of regexps - one for each level of trac heading
	compileRegexps()
}

func (converter *DefaultConverter) convertHeadings(in string) string {
	// iterate through all heading levels starting from longest (doing shortest first risks premature regexp matches)
	out := in
	for headingLevel := maxHeadingLevel - 1; headingLevel >= 0; headingLevel-- {
		headingRegexp := headingRegexps[headingLevel]
		out = headingRegexp.ReplaceAllStringFunc(out, func(match string) string {
			// turn trac delimiter into a markdown delimiter - e.g. '===' => '###'
			tracDelimiter := headingRegexp.ReplaceAllString(match, `$1`)
			markdownDelimiter := strings.Replace(tracDelimiter, "=", "#", -1)

			// extract text of heading
			headingText := headingRegexp.ReplaceAllString(match, `$2`)
			headingText = strings.Trim(headingText, " ")

			// examine any trailing text for presence of a Trac heading anchor
			anchor := ""
			trailingText := headingRegexp.ReplaceAllString(match, `$3`)
			if headingAnchorRegexp.MatchString(trailingText) {
				// if Trac anchor is the same as the "hyphenated" heading then this is the same as the implicit markdown heading anchor
				// so we don't need to embed an explicit anchor
				anchorName := headingAnchorRegexp.ReplaceAllString(trailingText, `$1`)
				hyphenatedHeading := strings.Replace(headingText, " ", "-", -1)
				if hyphenatedHeading != anchorName {
					// Trac anchor does not match markdown implicit anchor - the best we can do is insert a raw HTML anchor
					anchor = "<a name=\"" + anchorName + "\"></a>"
				}
			}

			return markdownDelimiter + " " + anchor + headingText
		})
	}

	return out
}
