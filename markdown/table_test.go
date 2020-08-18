package markdown_test

import "testing"

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

func TestSingleRowTable(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable := "||" + row1Cell1 + "||" + row1Cell2 + "||" + row1Cell3 + "||\n"

	// expect insertion of extra newline and for first row to be all headings
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n"

	conversion := converter.Convert(leadingText + "\n" + tracTable + trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+trailingText)
}

func TestMultiRowTable(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||" + row1Cell1 + "||" + row1Cell2 + "||" + row1Cell3 + "||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||\n"

	// expect insertion of extra newline and for first row to be all headings
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|\n"

	conversion := converter.Convert(leadingText + "\n" + tracTable + "\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}

func TestMultiRowTableWithAllHeadings(t *testing.T) {
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

	conversion := converter.Convert(leadingText + "\n" + tracTable + "\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}

func TestMultiRowTableWithSomeHeadings(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracTable :=
		"||=" + row1Cell1 + "=||" + row1Cell2 + "||=" + row1Cell3 + "=||\n" +
			"||" + row2Cell1 + "||" + row2Cell2 + "||" + row2Cell3 + "||\n" +
			"||" + row3Cell1 + "||" + row3Cell2 + "||" + row3Cell3 + "||\n"

	// expect insertion of extra newline and for first row to be all headings regardless of input
	markdownTable := "\n" +
		"|" + row1Cell1 + "|" + row1Cell2 + "|" + row1Cell3 + "|\n" +
		"|---|---|---|\n" +
		"|" + row2Cell1 + "|" + row2Cell2 + "|" + row2Cell3 + "|\n" +
		"|" + row3Cell1 + "|" + row3Cell2 + "|" + row3Cell3 + "|\n"

	conversion := converter.Convert(leadingText + "\n" + tracTable + "\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
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

	conversion := converter.Convert(leadingText + "\n" + tracTable + "\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownTable+"\n"+trailingText)
}
