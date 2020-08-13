package markdown

import "regexp"

var pageBreakRegexp = regexp.MustCompile(`\[\[[Bb][Rr]\]\]`)

func (converter *DefaultConverter) convertParagraphs(in string) string {
	return pageBreakRegexp.ReplaceAllString(in, "  \n")
}
