// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

import "regexp"

var blockQuoteRegexp = regexp.MustCompile(`(?m)^  ([[:alpha:]].*)$`)

func (converter *DefaultConverter) convertBlockQuotes(in string) string {
	return blockQuoteRegexp.ReplaceAllString(in, "> $1")
}
