// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import "testing"

const definition = "a definition"

func TestDefinitionList(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.WikiConvert(wikiPage, leadingText+"\n "+definition+"::"+trailingText)
	assertEquals(t, conversion, leadingText+"\n*"+definition+"*  \n"+trailingText)
}
