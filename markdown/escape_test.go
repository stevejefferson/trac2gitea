// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import "testing"

func TestEscape(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	escaped := "NotWikiLinkInTrac"

	conversion := converter.Convert(context, leadingText+"!"+escaped+trailingText)
	assertEquals(t, conversion, leadingText+escaped+trailingText)
}
