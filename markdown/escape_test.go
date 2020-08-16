package markdown_test

import "testing"

func TestEscape(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	escaped := "NotWikiLinkInTrac"

	conversion := converter.Convert(leadingText + "!" + escaped + trailingText)
	assertEquals(t, conversion, leadingText+escaped+trailingText)
}
