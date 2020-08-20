// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package markdown_test

import "testing"

const (
	headingText = "this is a heading"
)

func TestLevel1Heading(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n= " + headingText + " =\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n# "+headingText+"\n"+trailingText)
}

func TestLevel2Heading(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n== " + headingText + " ==\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n## "+headingText+"\n"+trailingText)
}

func TestLevel3Heading(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n=== " + headingText + " ===\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n### "+headingText+"\n"+trailingText)
}

func TestLevel4Heading(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n==== " + headingText + " ====\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n#### "+headingText+"\n"+trailingText)
}

func TestLevel5Heading(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n===== " + headingText + " =====\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n##### "+headingText+"\n"+trailingText)
}

func TestLevel6Heading(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n====== " + headingText + " ======\n" + trailingText)
	assertEquals(t, conversion, leadingText+"\n###### "+headingText+"\n"+trailingText)
}

func TestHeadingWithAnchor(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(leadingText + "\n==== " + headingText + " ==== #this-is-an-anchor\n" + trailingText)

	// ideally we should test for a warning being issued
	// unfortunately that would involve somehow intercepting/mocking our logging interface which is far from easy
	// so we will just have to content ourselves with testing that the conversion is unaffected by the anchor
	assertEquals(t, conversion, leadingText+"\n#### "+headingText+"\n"+trailingText)
}
