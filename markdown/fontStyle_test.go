package markdown_test

import "testing"

const (
	highlightedText = "some text which is highlighted in Trac"
)

func TestTripleSingleQuoteBold(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "'''" + highlightedText + "'''" + trailingText)
	assertEquals(t, conversion, leadingText+"**"+highlightedText+"**"+trailingText)
}
func TestDoubleSingleQuoteItalic(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "''" + highlightedText + "''" + trailingText)
	assertEquals(t, conversion, leadingText+"*"+highlightedText+"*"+trailingText)
}
func TestFiveSingleQuoteBoldItalic(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "'''''" + highlightedText + "'''''" + trailingText)
	assertEquals(t, conversion, leadingText+"**"+highlightedText+"**"+trailingText)
}
func TestDoubleAsteriskBold(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "**" + highlightedText + "**" + trailingText)
	assertEquals(t, conversion, leadingText+"**"+highlightedText+"**"+trailingText)
}
func TestDoubleSlashItalic(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "//" + highlightedText + "//" + trailingText)
	assertEquals(t, conversion, leadingText+"*"+highlightedText+"*"+trailingText)
}
func TestUnderline(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "__" + highlightedText + "__" + trailingText)
	assertEquals(t, conversion, leadingText+"*"+highlightedText+"*"+trailingText)
}
