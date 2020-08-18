package markdown

import (
	"fmt"
	"regexp"
)

var romanNumerals = []string{"i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix", "x", "xi", "xii", "xiii", "xiv", "xv", "xvi", "xvii", "xviii", "xix", "xx"}

// Numbered and bulleted Trac lists translate directly to markdown without translation .
// Therefore we only need to translate lettered lists ('a.', 'b.', ...) and roman-numbered lists ('i.', 'ii.', 'iv.' etc.).
// Due to the vaguaries of regular expressions we only handle lettered lists from 'a.' to 'h.'
// since 'i.' clashes with the roman-numbered case and is more likely to be the latter.
// For roman lists we do not police the actual numerals in the regexp - too painful.

// Regexp for lettered lists: $1=leading white space, $2=letter $3=trailing text
var letteredListRegexp = regexp.MustCompile(`(?m)^([[:blank:]]*)([a-h])\.([^\n]+)$`)

// Regexp for roman lists: $1=leading white space, $2=roman numerals $3=trailing text
var romanNumberedListRegexp = regexp.MustCompile(`(?m)^([[:blank:]]*)([ivx]+)\.([^\n]+)$`)

func (converter *DefaultConverter) convertLists(in string) string {
	out := in

	out = letteredListRegexp.ReplaceAllStringFunc(out, func(match string) string {
		leadingSpace := letteredListRegexp.ReplaceAllString(match, `$1`)
		letter := letteredListRegexp.ReplaceAllString(match, `$2`)
		trailingText := letteredListRegexp.ReplaceAllString(match, `$3`)
		letterNum := letter[0] - 'a' + 1 // 'a' => 1, 'b' => 2 etc
		return fmt.Sprintf("%s%d.%s", leadingSpace, letterNum, trailingText)
	})

	out = romanNumberedListRegexp.ReplaceAllStringFunc(out, func(match string) string {
		leadingSpace := romanNumberedListRegexp.ReplaceAllString(match, `$1`)
		roman := romanNumberedListRegexp.ReplaceAllString(match, `$2`)
		trailingText := romanNumberedListRegexp.ReplaceAllString(match, `$3`)
		for romanIndex := 0; romanIndex < len(romanNumerals); romanIndex++ {
			if romanNumerals[romanIndex] == roman {
				return fmt.Sprintf("%s%d.%s", leadingSpace, romanIndex+1, trailingText)
			}
		}

		return match
	})

	return out
}
