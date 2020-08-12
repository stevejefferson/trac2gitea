package markdown

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const maxHeadingLevel = 6

var headingRegexps [maxHeadingLevel]*regexp.Regexp

func compileRegexps() {
	tracHeadingDelimiter := ""
	for headingLevel := 0; headingLevel < maxHeadingLevel; headingLevel++ {
		tracHeadingDelimiter = tracHeadingDelimiter + "="

		// construct regexp for trac heading - this captures:
		// $1 = heading level delimiter (which we transform into the equivalent markdown delimiter)
		// $2 = heading text
		// $3 = optional Trac heading anchor (which we will typically have to warn about because markdown is not quote the same here)
		headingRegexpStr := `(?m)^(` + tracHeadingDelimiter + `) *([^=\n]+)` + tracHeadingDelimiter + `(?: +#([^[:space:]]+))?.*$`
		headingRegexps[headingLevel] = regexp.MustCompile(headingRegexpStr)
	}
}

func init() {
	// pre-compile array of regexps - one for each level of trac heading
	compileRegexps()
}

func (converter *Converter) convertHeadings(in string) string {
	// recurse through all heading levels starting from longest (doing shortest first risks premature regexp matches)
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

			// if any Trac anchor is given this must be the same as the heading after hyphenation
			// (markdown heading anchors are formed from the heading text and can't be arbitrary strings as in Trac)
			anchorName := headingRegexp.ReplaceAllString(match, `$3`)
			if anchorName != "" {
				hyphenatedHeading := strings.Replace(headingText, " ", "-", -1)
				if hyphenatedHeading != anchorName {
					fmt.Fprintf(os.Stderr, "Warning: anchor \"%s\" on trac heading \"%s\" cannot be used in markdown - hyphenate the heading text and use that to reference the heading\n",
						anchorName, headingText)
				}
			}

			return markdownDelimiter + " " + headingText
		})
	}

	return out
}
