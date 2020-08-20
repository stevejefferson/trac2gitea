// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package markdown_test

import "testing"

func TestSingleLineCodeBlock(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	code := "this is some code"

	conversion := converter.Convert(leadingText + "{{{" + code + "}}}" + trailingText)
	assertEquals(t, conversion, leadingText+"`"+code+"`"+trailingText)
}

func TestMultiLineCodeBlock(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	codeLine1 := "this is some code\n"
	codeLine2 := "this is more code\n"

	conversion := converter.Convert(
		leadingText + "\n" +
			"{{{#!trac-stuff\n" +
			codeLine1 +
			codeLine2 +
			"}}}\n" +
			trailingText)
	assertEquals(t, conversion,
		leadingText+"\n"+
			"```#!trac-stuff\n"+
			codeLine1+
			codeLine2+
			"```\n"+
			trailingText)
}

func TestNoConversionInsideCodeBlock(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	codeLine1 := "Website reference: http://www.example.com\n"
	codeLine2 := "[wiki:WikiPage, trac-style wiki link] followed by Trac-style //italics//\n"
	codeLine3 := "- bullet point\n"
	codeLine4 := "== Trac-style Subheading\n"

	conversion := converter.Convert(
		leadingText + "\n" +
			"{{{#!trac-stuff\n" +
			codeLine1 +
			codeLine2 +
			codeLine3 +
			codeLine4 +
			"}}}\n" +
			trailingText)
	assertEquals(t, conversion,
		leadingText+"\n"+
			"```#!trac-stuff\n"+
			codeLine1+
			codeLine2+
			codeLine3+
			codeLine4+
			"```\n"+
			trailingText)
}
