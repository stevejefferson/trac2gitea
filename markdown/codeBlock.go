package markdown

import (
	"regexp"
	"strings"
)

func (converter *Converter) convertCodeBlock(in string) string {
	// convert single line {{{...}}} to `...`
	out := in
	re := regexp.MustCompile("{{{([^\n]+?)}}}")
	out = re.ReplaceAllString(out, "`$1`")

	// convert multi-line {{{...}}} to tab-indented lines
	re = regexp.MustCompile("(?s){{{(.+?)}}}")
	out = re.ReplaceAllStringFunc(out, func(m string) string {
		lines := strings.Split(m, "\n")
		for i := range lines {
			l := lines[i]
			l = strings.Replace(l, "{{{", "", -1)
			l = strings.Replace(l, "}}}", "", -1)
			l = "\t" + l
			lines[i] = l
		}
		return strings.Join(lines, "\n")
	})

	return out
}
