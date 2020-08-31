// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import "testing"

const (
	listItem1 = "this is item 1"
	listItem2 = "this is item 2"
	listItem3 = "this is item 3"
)

func TestAsteriskBulletedLists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracList :=
		"* " + listItem1 + "\n" +
			"* " + listItem2 + "\n" +
			"* " + listItem3 + "\n"
	markdownList := tracList // asterisk bullets work in both trac and markdown

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracList+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownList+trailingText)
}

func TestHyphenBulletedLists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracList :=
		"- " + listItem1 + "\n" +
			"- " + listItem2 + "\n" +
			"- " + listItem3 + "\n"
	markdownList := tracList // hyphen bullets work in both trac and markdown

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracList+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownList+trailingText)
}

func TestNumberedBulletedLists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracList :=
		"1. " + listItem1 + "\n" +
			"2. " + listItem2 + "\n" +
			"3. " + listItem3 + "\n"
	markdownList := tracList // numbered bullets work in both trac and markdown

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracList+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownList+trailingText)
}

func TestLetteredBulletedLists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracList :=
		"a. " + listItem1 + "\n" +
			"b. " + listItem2 + "\n" +
			"f. " + listItem3 + "\n"
	markdownList :=
		"1. " + listItem1 + "\n" +
			"2. " + listItem2 + "\n" +
			"6. " + listItem3 + "\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracList+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownList+trailingText)
}
func TestRomanBulletedLists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracList :=
		"i. " + listItem1 + "\n" +
			"iv. " + listItem2 + "\n" +
			"xii. " + listItem3 + "\n"
	markdownList :=
		"1. " + listItem1 + "\n" +
			"4. " + listItem2 + "\n" +
			"12. " + listItem3 + "\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracList+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownList+trailingText)
}

func TestNestedLists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	tracList :=
		"1. " + listItem1 + "\n" +
			"  iv. " + listItem2 + "\n" +
			"    * " + listItem3 + "\n"
	markdownList :=
		"1. " + listItem1 + "\n" +
			"  4. " + listItem2 + "\n" +
			"    * " + listItem3 + "\n"

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n"+tracList+trailingText)
	assertEquals(t, conversion, leadingText+"\n"+markdownList+trailingText)
}
