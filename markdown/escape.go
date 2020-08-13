package markdown

import "regexp"

var escapeRegexp = regexp.MustCompile(`\!((?:[[:upper:]][[:lower:]]+){2,})`)

// convertEscapes converts Trac markdown escapes to Markdown
// - this must be run after convertLinks otherwise it will convert non-links into something that will be recognised as a link
func (converter *DefaultConverter) convertEscapes(in string) string {
	return escapeRegexp.ReplaceAllString(in, "$1")
}
