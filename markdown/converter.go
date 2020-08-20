// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package markdown

// Converter is the interface for Trac markdown to Gitea markdown conversions
type Converter interface {
	// Convert converts a string of Trac markdown to Gitea markdown
	Convert(in string) string
}
