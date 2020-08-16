package markdown_test

import "testing"

func TestUnlabelledAnchor(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "[=#name]" + trailingText)
	assertEquals(t, conversion, leadingText+"[](#name)"+trailingText)
}
func TestLabelledAnchor(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "[=#name anchor label]" + trailingText)
	assertEquals(t, conversion, leadingText+"[anchor label](#name)"+trailingText)
}
