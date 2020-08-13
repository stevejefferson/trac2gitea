package markdown

import "regexp"

var eolRegexp = regexp.MustCompile(`(?m)\r$`)

func (converter *DefaultConverter) convertEOL(in string) string {
	// Wiki lines within Trac database seem to have DOS-style `\r` terminated lines
	// - convert to Unix-style for consistency and to help some regexp matches
	return eolRegexp.ReplaceAllString(in, "")
}
