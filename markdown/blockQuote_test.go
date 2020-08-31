// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import "testing"

const (
	line1 = "this is line1\n"
	line2 = "this is line2\n"
)

func TestBlockQuote(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(
		context,
		leadingText+"\n"+
			"  "+line1+
			"  "+line2+
			trailingText)
	assertEquals(t, conversion,
		leadingText+"\n"+
			"> "+line1+
			"> "+line2+
			trailingText)
}
