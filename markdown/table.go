// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

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

func tracTableRowTransform(
	cells []string,
	cellSeparator string,
	actionIfCellIsHeader func(cell string) string,
	actionIfCellNotHeader func(cell string) string) string {
	result := cellSeparator
	for i := 0; i < len(cells); i++ {
		cell := cells[i]
		cellIsHeader := strings.HasPrefix(cell, "=") && strings.HasSuffix(cell, "=")
		var transformedCell string
		if cellIsHeader {
			transformedCell = actionIfCellIsHeader(cell)
		} else {
			transformedCell = actionIfCellNotHeader(cell)
		}

		result = result + transformedCell + cellSeparator
	}
	return result
}

func (converter *DefaultConverter) convertTables(in string) string {
	out := in

	// first pass: examine first row of table - this will need converting to be a header for the table to render in markdown
	// - note: output of this pass is a trac-format table, conversion to markdown happens in second pass
	out = tableFirstRowRegexp.ReplaceAllStringFunc(out, func(match string) string {
		prevLine := tableFirstRowRegexp.ReplaceAllString(match, `$1`)
		leadingWhitespace := tableFirstRowRegexp.ReplaceAllString(match, `$2`)
		rowContents := tableFirstRowRegexp.ReplaceAllString(match, `$3`)

		// if line prior to table is non-empty, insert a newline because markdown needs the table to be separate for preceeding content
		if prevLine != "" {
			prevLine = prevLine + "\n"
		}

		// determine whether any of the cells in the first row of the trac table are headers
		cells := strings.Split(rowContents, "||")
		haveAHeaderCell := false
		tracTableRowTransform(cells, "||",
			func(cell string) string {
				haveAHeaderCell = true
				return cell
			},
			func(cell string) string {
				return cell
			})

		var header string
		if haveAHeaderCell {
			// first row contains at least one header: turn all cells in this row into header cells
			header = tracTableRowTransform(cells, "||",
				func(cell string) string {
					return cell
				},
				func(cell string) string {
					return "=" + cell + "=" // make non-header into a header
				})
		} else {
			// first row does not contain any headers: prepend a blank header row to the table
			header = tracTableRowTransform(cells, "||",
				func(cell string) string {
					return "= ="
				},
				func(cell string) string {
					return "= ="
				})

			header = header + "\n||" + rowContents + "||"
		}

		return prevLine + "\n" + leadingWhitespace + header + "\n"
	})

	// second pass: convert all trac table rows to markdown
	out = tableRowRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// split row into cells
		leadingWhitespace := tableRowRegexp.ReplaceAllString(match, `$1`)
		rowContents := tableRowRegexp.ReplaceAllString(match, `$2`)
		cells := strings.Split(rowContents, "||")

		// build up markdown-format cell row plus an optional header row if any Trac '||=...=||' cells are found
		haveAHeaderCell := false
		cellRow := tracTableRowTransform(cells, "|",
			func(cell string) string {
				haveAHeaderCell = true
				return cell[1 : len(cell)-1] // strip trac '=' delimiters off cell
			},
			func(cell string) string {
				return cell
			})
		headerRow := tracTableRowTransform(cells, "|",
			func(cell string) string {
				return "---"
			},
			func(cell string) string {
				return ""
			})

		result := cellRow
		if haveAHeaderCell {
			result = result + "\n" + headerRow
		}

		return leadingWhitespace + result + "\n"
	})

	return out
}
