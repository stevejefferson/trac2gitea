// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package markdown

import (
	"regexp"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
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
		//      (we have to warn about because markdown heading anchors are not the same)
		headingRegexpStr := `(?m)^(` + tracHeadingDelimiter + `) *([^=\n]+)` + tracHeadingDelimiter + `(.*)$`
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
			trailingText := headingRegexp.ReplaceAllString(match, `$3`)
			if headingAnchorRegexp.MatchString(trailingText) {
				// any Trac anchor must be the same as the heading after hyphenation
				// (markdown heading anchors are formed from the heading text and can't be arbitrary strings as in Trac)
				anchorName := headingAnchorRegexp.ReplaceAllString(trailingText, `$1`)
				hyphenatedHeading := strings.Replace(headingText, " ", "-", -1)
				if hyphenatedHeading != anchorName {
					log.Warnf("anchor \"%s\" on trac heading \"%s\" cannot be used in markdown - hyphenate the heading text and use that to reference the heading\n",
						anchorName, headingText)
				}
			}

			return markdownDelimiter + " " + headingText
		})
	}

	return out
}
