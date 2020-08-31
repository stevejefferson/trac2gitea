// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown

// ConversionContext provides the context for a conversion - either conversion of a ticket (comment) or a wiki page
type ConversionContext struct {
	TicketID int64
	WikiPage string
}

// Converter is the interface for Trac markdown to Gitea markdown conversions
type Converter interface {
	// Convert converts a Trac markdown string to Gitea markdown
	Convert(context *ConversionContext, in string) string
}
