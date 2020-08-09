package markdown

import (
	"regexp"
)

const maxHeadingLevel = 6

var headingRegexps [maxHeadingLevel]*regexp.Regexp
var headingReplaceStrs [maxHeadingLevel]string

func compileRegexps(headingLevel int, delimiter string) {
	if headingLevel >= maxHeadingLevel {
		return
	}

	delimiterChar := delimiter[0:1]
	headingRegexpStr := `(?m)^` + delimiter + `([^` + delimiterChar + `]+)` + delimiter + `.*$`
	headingRegexps[headingLevel] = regexp.MustCompile(headingRegexpStr)
	compileRegexps(headingLevel+1, delimiter+delimiterChar)
}

func createHeadingReplaceStrs(headingLevel int, delimiter string) {
	if headingLevel >= maxHeadingLevel {
		return
	}

	delimiterChar := delimiter[0:1]
	headingReplaceStr := delimiter + `$1`
	headingReplaceStrs[headingLevel] = headingReplaceStr
	createHeadingReplaceStrs(headingLevel+1, delimiter+delimiterChar)
}

func init() {
	// pre-compile array of regexps - one for each level of trac heading
	compileRegexps(0, `=`)

	// generate markdown replacement strings
	createHeadingReplaceStrs(0, `#`)
}

func (converter *Converter) convertHeading(in string) string {
	// recurse through all heading levels starting from longest (doing shortest first risks premature regexp matches)
	out := in
	for headingLevel := maxHeadingLevel - 1; headingLevel >= 0; headingLevel-- {
		out = headingRegexps[headingLevel].ReplaceAllString(out, headingReplaceStrs[headingLevel])
	}

	return out
}
