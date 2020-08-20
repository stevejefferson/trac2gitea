// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

import "regexp"

var definitionListRegexp = regexp.MustCompile(`(?m)^ ([^:]+)\:\:`)

func (converter *DefaultConverter) convertDefinitionLists(in string) string {
	return definitionListRegexp.ReplaceAllString(in, "*$1*  \n")
}
