// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import "testing"

func TestUpperCasePageBreak(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(context, leadingText+"[[BR]]"+trailingText)
	assertEquals(t, conversion, leadingText+"<br>"+trailingText)
}

func TestLowerCasePageBreak(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	conversion := converter.Convert(context, leadingText+"[[br]]"+trailingText)
	assertEquals(t, conversion, leadingText+"<br>"+trailingText)
}
