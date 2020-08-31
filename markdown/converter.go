// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

// Converter is the interface for Trac markdown to Gitea markdown conversions
type Converter interface {
	// TicketConvert converts a comment/description string associated with a Trac ticket to Gitea markdown
	TicketConvert(ticketID int64, in string) string

	// WikiConvert converts a comment/description string associated with a Trac wiki page to Gitea markdown
	WikiConvert(wikiPage string, in string) string
}
