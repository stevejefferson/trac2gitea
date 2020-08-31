// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import (
	"testing"
)

const (
	row1Cell1 = "this is row 1, cell 1"
	row1Cell2 = "this is row 1, cell 2"
	row1Cell3 = "this is row 1, cell 3"

	row2Cell1 = "this is row 2, cell 1"
	row2Cell2 = "this is row 2, cell 2"
	row2Cell3 = "this is row 2, cell 3"

	row3Cell1 = "this is row 3, cell 1"
	row3Cell2 = "this is row 3, cell 2"
	row3Cell3 = "this is row 3, cell 3"
)

func TestSingleNonHeaderRowTableHasHeaderRowPrepended(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable := "||" + row1Cell1 + "||" + row1Cell2 + "||" + row1Cell3 + "||\n"
	markdownTable := "\n" +
		"| | | |\n" +
		"|---|---|---|\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n"
	conversion := converter.WikiConvert(wikiPage, leadingText+"\n\n"+tracTable+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+trailingText)
}

func TestSingleNonHeaderRowTableAtStartOfTextHasHeaderRowPrepended(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable := "||" + row1Cell1 + "||" + row1Cell2 + "||" + row1Cell3 + "||\n"
	markdownTable :=
		"| | | |\n" +
			"|---|---|---|\n" +
			"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n"
	conversion := converter.WikiConvert(wikiPage, tracTable+trailingText)
	assertEquals(t, conversion, markdownTable+trailingText)
}

func TestSinglePartialHeaderRowTableBecomesAllHeaderRow(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable := "||" + row1Cell1 + "||=" + row1Cell2 + "=||" + row1Cell3 + "||\n"
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n"
	conversion := converter.WikiConvert(wikiPage, leadingText+"\n\n"+tracTable+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+trailingText)
}

func TestSingleAllHeaderRowTableRemainsAllHeaderRow(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable := "||=" + row1Cell1 + "=||=" + row1Cell2 + "=||=" + row1Cell3 + "=||\n"
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n"
	conversion := converter.WikiConvert(wikiPage, leadingText+"\n\n"+tracTable+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+trailingText)
}

func TestBlankLineInsertedBetweenPrevLineAndTable(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable := "||=" + row1Cell1 + "=||=" + row1Cell2 + "=||=" + row1Cell3 + "=||\n"
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n"

	// note omission of "\n" in text to convert compared to prev test
	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+trailingText)
}

func TestMultiRowTableWithNoHeader(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||" + row1Cell1 + "||" + row1Cell2 + "||" + row1Cell3 + "||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||\n"

	markdownTable := "\n" +
		"| | | |\n" +
		"|---|---|---|\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable+"\n"+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}

func TestMultiRowTableWithPartialHeader(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||" + row1Cell1 + "||" + row1Cell2 + "||=" + row1Cell3 + "=||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||\n"

	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable+"\n"+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}

func TestMultiRowTableWithAllHeaders(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||=" + row1Cell1 + "=||=" + row1Cell2 + "=||=" + row1Cell3 + "=||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||\n"

	// expect insertion of extra newline and for first row to (still) be all headings
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable+"\n"+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}

func TestTableAtStartOfText(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||=" + row1Cell1 + "=||=" + row1Cell2 + "=||=" + row1Cell3 + "=||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||"

	// expect insertion of extra newline and for first row to (still) be all headings
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable)
}

func TestTableAtEndOfText(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||=" + row1Cell1 + "=||=" + row1Cell2 + "=||=" + row1Cell3 + "=||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||"

	// expect insertion of extra newline and for first row to (still) be all headings
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable)
}

func TestMultiRowTableWithCrazyHeadings(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||=" + row1Cell1 + "=||" + row1Cell2 + "||" + row1Cell3 + "||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||=" + row2Cell3 + "=||\n" +
			"||" + row3Cell1 + "||=" + row3Cell2 + "=||" + row3Cell3 + "||\n"

	// expect insertion of extra newline and for first row to be all headings regardless of input
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|||---|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|\n" +
		"||---||\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracTable+"\n"+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}
