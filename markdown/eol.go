package markdown

import "regexp"

func (converter *Converter) convertEOL(in string) string {
	// Wiki lines within Trac database seem to have DOS-style `\r` terminated lines
	// - convert to Unix-style `\n` otherwise '$'-terminated regexps don't seem to work
	regexStr := `\r`
	regex := regexp.MustCompile(regexStr)
	return regex.ReplaceAllString(in, "")
}
