package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

var romanNumerals = []string{"i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix", "x", "xi", "xii", "xiii", "xiv", "xv", "xvi", "xvii", "xviii", "xix", "xx"}

// Numbered and bulleted Trac lists translate directly to markdown without translation .
// Therefore we only need to translate lettered lists ('a.', 'b.', ...) amd roman-numbered lists ('i.', 'ii.', 'iv.' etc.).
// Due to the vaguaries of regular expressions be only handle lettered lists from 'a.' to 'h.'
// since 'i.' clashes with the roman-numbered case and is more likely to be the latter.
// For roman lists we do not police the actual numerals in the regexp - too painful.
var letteredListRegexp = regexp.MustCompile(`[a-h]\. [^\n]+`)
var romanNumberedListRegexp = regexp.MustCompile(`[ivx]+\. ([^\n]+)`)

func (converter *DefaultConverter) convertLists(in string) string {
	out := in

	out = letteredListRegexp.ReplaceAllStringFunc(out, func(match string) string {
		letterNum := match[0] - 'a' + 1 // 'a' => 1, 'b' => 2 etc
		return fmt.Sprintf("%d%s", letterNum, match[1:])
	})

	out = romanNumberedListRegexp.ReplaceAllStringFunc(out, func(match string) string {
		dotPos := strings.Index(match, ".")
		romanNumeral := match[0:dotPos]
		for i := 0; i < len(romanNumerals); i++ {
			if romanNumerals[i] == romanNumeral {
				return fmt.Sprintf("%d%s", i+1, match[dotPos:])
			}
		}

		return match
	})

	return out
}
