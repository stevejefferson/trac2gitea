package markdown_test

import "testing"

func TestUpperCasePageBreak(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "[[BR]]" + trailingText)
	assertEquals(t, conversion, leadingText+"<br>"+trailingText)
}

func TestLowerCasePageBreak(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "[[br]]" + trailingText)
	assertEquals(t, conversion, leadingText+"<br>"+trailingText)
}
