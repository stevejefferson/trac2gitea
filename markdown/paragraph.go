package markdown

import "regexp"

var pageBreakRegexp = regexp.MustCompile(`\[\[[Bb][Rr]\]\]`)

func (converter *Converter) convertParagraphs(in string) string {
	return pageBreakRegexp.ReplaceAllString(in, "  \n")
}
