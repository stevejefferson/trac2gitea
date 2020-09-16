// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

// Milestone describes a Trac milestone.
type Milestone struct {
	Name        string
	Description string
	Due         int64
	Completed   int64
}

const (
	// TicketStatusClosed indicates a closed Trac ticket
	TicketStatusClosed string = "closed"
)

// Ticket describes a Trac milestone.
type Ticket struct {
	TicketID       int64
	Summary        string
	Description    string
	Owner          string
	Reporter       string
	MilestoneName  string
	ComponentName  string
	PriorityName   string
	ResolutionName string
	SeverityName   string
	TypeName       string
	VersionName    string
	Status         string
	Created        int64
	Updated        int64
}

// TicketChangeType enumerates the types of ticket change we handle.
type TicketChangeType string

const (
	// TicketCommentChange denotes a ticket comment change.
	TicketCommentChange TicketChangeType = "comment"

	// TicketComponentChange denotes a ticket component change.
	TicketComponentChange TicketChangeType = "component"

	// TicketMilestoneChange denotes a ticket milestone change.
	TicketMilestoneChange TicketChangeType = "milestone"

	// TicketOwnerChange denotes a ticket ownership change.
	TicketOwnerChange TicketChangeType = "owner"

	// TicketPriorityChange denotes a ticket resolution change.
	TicketPriorityChange TicketChangeType = "priority"

	// TicketResolutionChange denotes a ticket resolution change.
	TicketResolutionChange TicketChangeType = "resolution"

	// TicketSeverityChange denotes a ticket severity change.
	TicketSeverityChange TicketChangeType = "severity"

	// TicketStatusChange denotes a ticket status change.
	TicketStatusChange TicketChangeType = "status"

	// TicketSummaryChange denotes a ticket summary change.
	TicketSummaryChange TicketChangeType = "summary"

	// TicketTypeChange denotes a ticket type change.
	TicketTypeChange TicketChangeType = "type"

	// TicketVersionChange denotes a ticket type change.
	TicketVersionChange TicketChangeType = "version"
)

// TicketChange describes a change to a Trac ticket.
type TicketChange struct {
	TicketID   int64
	ChangeType TicketChangeType
	Author     string
	OldValue   string
	NewValue   string
	Time       int64
}

// TicketAttachment describes an attachment to a Trac ticket.
type TicketAttachment struct {
	TicketID    int64
	Time        int64
	Size        int64
	Author      string
	FileName    string
	Description string
}

// WikiPage describes a Trac wiki page.
type WikiPage struct {
	Name       string
	Text       string
	Author     string
	Comment    string
	Version    int64
	UpdateTime int64
}

// WikiAttachment describes an attachment to a Trac wiki page.
type WikiAttachment struct {
	PageName string
	FileName string
}

// Accessor is the interface through which we access all Trac data.
type Accessor interface {
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
	GetMilestones(handlerFn func(milestone *Milestone) error) error

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
	GetTickets(handlerFn func(ticket *Ticket) error) error

	/*
	 * Ticket Changes
	 */
	// GetTicketChanges retrieves all changes on a given Trac ticket in ascending time order, passing data from each one to the provided "handler" function.
	GetTicketChanges(ticketID int64, handlerFn func(change *TicketChange) error) error

	// GetTicketCommentTime retrieves the timestamp for a given comment for a given Trac ticket
	GetTicketCommentTime(ticketID int64, changeNum int64) (int64, error)

	/*
	 * Ticket Attachments
	 */
	// GetTicketAttachmentPath retrieves the path to a named attachment to a Trac ticket.
	GetTicketAttachmentPath(attachment *TicketAttachment) string

	// GetTicketAttachments retrieves all attachments for a given Trac ticket, passing data from each one to the provided "handler" function.
	GetTicketAttachments(ticketID int64, handlerFn func(attachment *TicketAttachment) error) error

	/*
	 * Types
	 */
	// GetTypeNames retrieves all type names used in Trac tickets, passing each one to the provided "handler" function.
	GetTypeNames(handlerFn func(typeName string) error) error

	/*
	 * Users
	 */
	// GetUserNames retrieves the names of all users mentioned in Trac tickets, wiki pages etc., passing each one to the provided "handler" function.
	GetUserNames(handlerFn func(userName string) error) error

	/*
	 * Versions
	 */
	// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
	GetVersionNames(handlerFn func(version string) error) error

	/*
	 * Wiki
	 */
	// GetWikiPages retrieves all Trac wiki pages, passing data from each one to the provided "handler" function.
	GetWikiPages(handlerFn func(page *WikiPage) error) error

	// GetWikiAttachmentPath retrieves the path to a named attachment to a Trac wiki page.
	GetWikiAttachmentPath(attachment *WikiAttachment) string

	// GetWikiAttachments retrieves all Trac wiki page attachments, passing data from each one to the provided "handler" function.
	GetWikiAttachments(handlerFn func(attachment *WikiAttachment) error) error

	// IsPredefinedPage returns true if the provided page name is one of Trac's predefined ones - by default we ignore these
	IsPredefinedPage(pageName string) bool
}
