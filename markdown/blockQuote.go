package markdown

import "regexp"

var blockQuoteRegexp = regexp.MustCompile(`(?m)^  ([[:alpha:]].*)$`)

func (converter *Converter) convertBlockQuotes(in string) string {
	return blockQuoteRegexp.ReplaceAllString(in, "> $1")
}
