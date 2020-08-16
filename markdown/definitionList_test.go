package markdown_test

import "testing"

const definition = "a definition"

func TestDefinitionList(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n " + definition + "::" + trailingText)
	assertEquals(t, conversion, leadingText+"\n*"+definition+"*  \n"+trailingText)
}
