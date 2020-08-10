package markdown

import "regexp"

var anchorRegexp = regexp.MustCompile(`(?m)\[=#([^\]\n ]+?)( +[^\]\n]+)?\]`)

// convertAnchors converts Trac '[=#name...]' anchors
// additionally Trac supports anchors on headings - these are dealt with in the heading conversion
func (converter *Converter) convertAnchors(in string) string {
	out := in
	out = anchorRegexp.ReplaceAllStringFunc(out, func(match string) string {
		anchorName := anchorRegexp.ReplaceAllString(match, `$1`)
		anchorLabel := anchorRegexp.ReplaceAllString(match, `$2`)
		return "[" + anchorLabel + "](#" + anchorName + ")"
	})

	return out
}
