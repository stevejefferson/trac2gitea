// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

import (
	"regexp"
	"strings"
)

var (
	// regexp matching the first row of a trac table at the start of text
	// $1=leading whitespace prior to first row, $2=first row starting from after first '||'
	tableFirstRowStartOfTextRegexp = regexp.MustCompile(`(?m)\A([[:blank:]]*)\|\|([^\n]*)$`)

	// regexp matching the first row of a trac table:
	// $1=previous (non-table) line, $2=leading whitespace prior to first row, $3=first row starting from after first '||'
	tableFirstRowRegexp = regexp.MustCompile(`(?m)^([[:blank:]]*(?:[^\|\n][^\n]*)?)$\n([[:blank:]]*)\|\|([^\n]*)$`)

	// regexp matching entire trac table row including any intermediate cell separators:
	// $1=leading whitespace, $2=row contents starting from after first '||'
	tableRowRegexp = regexp.MustCompile(`(?m)([[:blank:]]*)\|\|(.*)$`)
)

// tracTableRowTransform is a general transformation function for identifying and processing cells in a row from a trac table
func tracTableRowTransform(
	tracRow string,
	cellSeparator string,
	actionIfCellIsHeader func(cell string) string,
	actionIfCellNotHeader func(cell string) string) string {

	// split table row into cells
	// - input is always a trac-format row so always separated by '||'
	// - remember that last one is the text between the terminating '||' and the end-of-line so should be skipped
	cells := strings.Split(tracRow, "||")

	result := cellSeparator
	for i := 0; i < len(cells)-1; i++ {
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

// makeFirstRowAHeader makes the first row of a trac table into a header either by making any non-header cells into header cells
// or by inserting a blank header row if no cells are headers.
// The result of this function is in trac table format.
func makeFirstRowAHeader(leadingWhitespace string, tracRow string) string {
	// determine whether any of the cells in the row are headers
	haveAHeaderCell := false
	tracTableRowTransform(tracRow, "||",
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
		header = tracTableRowTransform(tracRow, "||",
			func(cell string) string {
				return cell
			},
			func(cell string) string {
				return "=" + cell + "=" // make non-header into a header
			})
	} else {
		// first row does not contain any headers: prepend a blank header row to the table
		header = tracTableRowTransform(tracRow, "||",
			func(cell string) string {
				return "= ="
			},
			func(cell string) string {
				return "= ="
			})

		header = header + "\n" + leadingWhitespace + "||" + tracRow
	}

	return leadingWhitespace + header
}

func tracRowToMarkdown(leadingWhitespace string, tracRow string) string {
	// build up markdown-format cell row plus an optional header row if any Trac '||=...=||' cells are found
	haveAHeaderCell := false
	cellRow := tracTableRowTransform(tracRow, "|",
		func(cell string) string {
			haveAHeaderCell = true
			return cell[1 : len(cell)-1] // strip trac '=' delimiters off cell
		},
		func(cell string) string {
			return cell
		})
	result := cellRow

	if haveAHeaderCell {
		headerRow := tracTableRowTransform(tracRow, "|",
			func(cell string) string {
				return "---"
			},
			func(cell string) string {
				return ""
			})

		result = result + "\n" + headerRow
	}

	result = leadingWhitespace + result
	return result
}

func (converter *DefaultConverter) convertTables(in string) string {
	out := in

	// First pass: examine first row of table - this will need converting to be a header for the table to render in markdown.
	// Due to the way we phrase our regexps we need to do this in two parts - one where the first row is at the start of text
	// and one where it is identified by a preceeding non-table line.
	// Note: in both cases, the output of this pass is a trac-format table, conversion to markdown happens in second pass
	out = tableFirstRowStartOfTextRegexp.ReplaceAllStringFunc(out, func(match string) string {
		leadingWhitespace := tableFirstRowStartOfTextRegexp.ReplaceAllString(match, `$1`)
		rowContents := tableFirstRowStartOfTextRegexp.ReplaceAllString(match, `$2`)

		rowWithHeader := makeFirstRowAHeader(leadingWhitespace, rowContents)
		return rowWithHeader
	})

	out = tableFirstRowRegexp.ReplaceAllStringFunc(out, func(match string) string {
		prevLine := tableFirstRowRegexp.ReplaceAllString(match, `$1`)
		leadingWhitespace := tableFirstRowRegexp.ReplaceAllString(match, `$2`)
		rowContents := tableFirstRowRegexp.ReplaceAllString(match, `$3`)

		// if line prior to table is non-empty, insert a newline because markdown needs the table to be separate for preceeding content
		if prevLine != "" {
			prevLine = prevLine + "\n"
		}

		rowWithHeader := makeFirstRowAHeader(leadingWhitespace, rowContents)
		result := prevLine + "\n" + rowWithHeader
		return result
	})

	// second pass: convert all trac table rows to markdown
	out = tableRowRegexp.ReplaceAllStringFunc(out, func(match string) string {
		// split row into cells
		leadingWhitespace := tableRowRegexp.ReplaceAllString(match, `$1`)
		rowContents := tableRowRegexp.ReplaceAllString(match, `$2`)
		markdownRow := tracRowToMarkdown(leadingWhitespace, rowContents)
		return markdownRow
	})

	return out
}
