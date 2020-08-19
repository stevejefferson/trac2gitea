package markdown

import (
	"regexp"
	"strings"
)

var (
	// regexp matching the first row of a trac table: $1=previous (non-table) line, $2=leading whitespace prior to first row, $3=first row
	tableFirstRowRegexp = regexp.MustCompile(`(?m)^([[:blank:]]*(?:[^\|\n][^\n]*)?)\n([[:blank:]]*)\|\|([^\n]*)\|\|.*\n`)

	// regexp matching entire trac table row including any intermediate cell separators: $1=leading whitespace, $2=row contents
	tableRowRegexp = regexp.MustCompile(`(?m)([[:blank:]]*)\|\|([^\n]*)\|\|[^\n]*\n`)
)

func (converter *DefaultConverter) convertTables(in string) string {
	out := in

	out = tableFirstRowRegexp.ReplaceAllStringFunc(out, func(match string) string {
		prevLine := tableFirstRowRegexp.ReplaceAllString(match, `$1`)
		leadingWhitespace := tableFirstRowRegexp.ReplaceAllString(match, `$2`)
		rowContents := tableFirstRowRegexp.ReplaceAllString(match, `$3`)

		// if line prior to table is non-empty, insert a newline because markdown needs the table to be separate
		result := prevLine
		if prevLine != "" {
			result = result + "\n"
		}

		result = result + "\n" + leadingWhitespace + "||"

		// ensure all cells in first row are Trac header cells (next pass will convert them to markdown)
		// - unless first row is a header row, our table won't display
		cells := strings.Split(rowContents, "||")
		for i := 0; i < len(cells); i++ {
			cell := cells[i]
			cellIsHeader := strings.HasPrefix(cell, "=") && strings.HasSuffix(cell, "=")
			if !cellIsHeader {
				cell = "=" + cell + "="
			}
			result = result + cell + "||"
		}
		return result + "\n"
	})

	out = tableRowRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// split row into cells
		leadingWhitespace := tableRowRegexp.ReplaceAllString(match, `$1`)
		rowContents := tableRowRegexp.ReplaceAllString(match, `$2`)
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

		return leadingWhitespace + result + "\n"
	})

	return out
}
