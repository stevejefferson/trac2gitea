package markdown

import (
	"regexp"
	"strings"
)

var singleLineCodeBlockRegexp = regexp.MustCompile(`{{{([^\n]+?)}}}`)
var multiLineCodeBlockRegexp = regexp.MustCompile(`(?m)^{{{(?s)(.+?)^}}}`)
var nonCodeBlockRegexp = regexp.MustCompile(`(?m)(?:}}}$|\A)(?s)(.+?)(?:^{{{|\z)`)

func (converter *Converter) convertCodeBlocks(in string) string {
	// convert single line {{{...}}} to `...`
	out := singleLineCodeBlockRegexp.ReplaceAllString(in, "`$1`")

	// convert multi-line {{{...}}} to tab-indented lines
	// - we leave in place any Trac '#!...' sequences following the opening '{{{'
	//   since we have no easy way of dealing with these and they are best left in place
	//   as a reminder to review them in the Gitea world
	out = multiLineCodeBlockRegexp.ReplaceAllStringFunc(out, func(match string) string {
		text := multiLineCodeBlockRegexp.ReplaceAllString(match, `$1`)

		lines := strings.Split(text, "\n")
		for i := range lines {
			lines[i] = "\t" + lines[i]
		}
		return strings.Join(lines, "\n")
	})

	return out
}

func (converter *Converter) convertNonCodeBlocks(in string, convertFn func(string) string) string {
	return nonCodeBlockRegexp.ReplaceAllStringFunc(in, convertFn)
}
