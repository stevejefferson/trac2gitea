package markdown

import (
	"regexp"
	"strings"
)

var tableRegexp = regexp.MustCompile(`(?m)\|\|(.*)\|\|`) // matches entire table row including any intermediate cell separators

func (converter *Converter) convertTables(in string) string {
	out := in
	out = tableRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// split row into cells
		rowContents := tableRegexp.ReplaceAllString(match, `$1`)
		cells := strings.Split(rowContents, "||")

		// build up markdown-format cell row plus an optional header row if any Trac '||=...=||' cells are found
		cellRow := "|"
		headerRow := "|"
		haveAHeaderCell := false
		for i := 0; i < len(cells); i++ {
			cell := cells[i]
			cellIsHeader := strings.HasPrefix(cell, "=") && strings.HasSuffix(cell, "=")

			var cellText = cell
			var headerText = ""
			if cellIsHeader {
				haveAHeaderCell = true
				cellText = cell[1 : len(cell)-1]
				headerText = "---"
			}

			cellRow = cellRow + cellText + "|"
			headerRow = headerRow + headerText + "|"
		}

		result := cellRow
		if haveAHeaderCell {
			result = result + "\n" + headerRow
		}

		return result
	})

	return out
}
