package markdown

import (
	"regexp"
	"strings"
)

var singleLineCodeBlockRegexp = regexp.MustCompile(`{{{([^\n]+?)}}}`)
var multiLineCodeBlockRegexp = regexp.MustCompile(`(?s){{{(.+?)}}}`)

func (converter *Converter) convertCodeBlock(in string) string {
	// convert single line {{{...}}} to `...`
	out := singleLineCodeBlockRegexp.ReplaceAllString(in, "`$1`")

	// convert multi-line {{{...}}} to tab-indented lines
	out = multiLineCodeBlockRegexp.ReplaceAllStringFunc(out, func(match string) string {
		lines := strings.Split(match, "\n")
		for i := range lines {
			line := lines[i]
			line = strings.Replace(line, "{{{", "", -1)
			line = strings.Replace(line, "}}}", "", -1)
			line = "\t" + line
			lines[i] = line
		}
		return strings.Join(lines, "\n")
	})

	return out
}
