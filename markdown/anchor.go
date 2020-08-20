// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

import "regexp"

// regexp for a Trac anchor: $1=anchor $2=anchor text
var anchorRegexp = regexp.MustCompile(`(?m)\[=#([[:alnum:]?/:@\-._\~!$&'()*+,;=]+)(?: +([^\]\n]+))?\]`)

// convertAnchors converts Trac '[=#name...]' anchors
// additionally Trac supports anchors on headings - these are dealt with in the heading conversion
func (converter *DefaultConverter) convertAnchors(in string) string {
	out := in
	out = anchorRegexp.ReplaceAllStringFunc(out, func(match string) string {
		anchorName := anchorRegexp.ReplaceAllString(match, `$1`)
		anchorLabel := anchorRegexp.ReplaceAllString(match, `$2`)
		return "[" + anchorLabel + "](#" + anchorName + ")"
	})

	return out
}
