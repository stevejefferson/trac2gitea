// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

// Accessor is the interface through which we access all Trac data.
type Accessor interface {
	/*
	 * Attachments
	 */
	// GetTicketAttachmentPath retrieves the path to a named attachment to a Trac ticket.
	GetTicketAttachmentPath(ticketID int64, attachmentName string) string

	// GetWikiAttachmentPath retrieves the path to a named attachment to a Trac wiki page.
	GetWikiAttachmentPath(wikiPage string, attachmentName string) string

	// GetAttachments retrieves all attachments for a given Trac ticket, passing data from each one to the provided "handler" function.
	GetAttachments(ticketID int64,
		handlerFn func(ticketID int64, time int64, size int64, author string, filename string, description string) error) error

	/*
	 * Comments
	 */
	// GetComments retrieves all comments on a given Trac ticket, passing data from each one to the provided "handler" function.
	GetComments(ticketID int64,
		handlerFn func(ticketID int64, time int64, author string, comment string) error) error

	// GetCommentString retrieves a given comment string for a given Trac ticket
	GetCommentString(ticketID int64, commentNum int64) (string, error)

	/*
	 * Components
	 */
	// GetComponentNames retrieves all Trac component names, passing each one to the provided "handler" function.
	GetComponentNames(handlerFn func(cmptName string) error) error

	/*
	 * Configuration
	 */
	// GetStringConfig retrieves a value from the Trac config as a string.
	GetStringConfig(sectionName string, configName string) string

	/*
	 * Milestones
	 */
	// GetMilestones retrieves all Trac milestones, passing data from each one to the provided "handler" function.
	GetMilestones(handlerFn func(name string, description string, due int64, completed int64) error) error

	/*
	 * Paths
	 */
	// GetFullPath retrieves the absolute path of a path relative to the root of the Trac installation.
	GetFullPath(element ...string) string

	/*
	 * Priorities
	 */
	// GetPriorityNames retrieves all priority names used in Trac tickets, passing each one to the provided "handler" function.
	GetPriorityNames(handlerFn func(priorityName string) error) error

	/*
	 * Resolutions
	 */
	// GetResolutionNames retrieves all resolution names used in Trac tickets, passing each one to the provided "handler" function.
	GetResolutionNames(handlerFn func(resolution string) error) error

	/*
	 * Severities
	 */
	// GetSeverityNames retrieves all severity names used in Trac tickets, passing each one to the provided "handler" function.
	GetSeverityNames(handlerFn func(severityName string) error) error

	/*
	 * Tickets
	 */
	// GetTickets retrieves all Trac tickets, passing data from each one to the provided "handler" function.
	GetTickets(handlerFn func(
		ticketID int64, ticketType string, created int64,
		component string, severity string, priority string,
		owner string, reporter string, version string,
		milestone string, status string, resolution string,
		summary string, description string) error) error

	/*
	 * Types
	 */
	// GetTypeNames retrieves all type names used in Trac tickets, passing each one to the provided "handler" function.
	GetTypeNames(handlerFn func(typeName string) error) error

	/*
	 * Users
	 */
	// GetUserMap returns a blank user mapping mapping for every user name found in Trac database fields to be converted
	GetUserMap() (map[string]string, error)

	/*
	 * Versions
	 */
	// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
	GetVersionNames(handlerFn func(version string) error) error

	/*
	 * Wiki
	 */
	// GetWikiPages retrieves all Trac wiki pages, passing data from each one to the provided "handler" function.
	GetWikiPages(handlerFn func(pageName string, pageText string, author string, comment string, version int64, updateTime int64) error) error

	// GetWikiAttachments retrieves all Trac wiki page attachments, passing data from each one to the provided "handler" function.
	GetWikiAttachments(handlerFn func(wikiPage string, filename string) error) error

	// IsPredefinedPage returns true if the provided page name is one of Trac's predefined ones - by default we ignore these
	IsPredefinedPage(pageName string) bool
}
