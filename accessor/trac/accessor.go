// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

// Label describes a Trac "label" - a generalisation of a component, priority, resolution, severity, type and version
type Label struct {
	Name        string
	Description string
}

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

	// TicketStatusReopened indicates a reopened Trac ticket
	TicketStatusReopened string = "reopened"
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

// NullID id used for Trac lookup failures
// The Trac schema does not seem to use foreign key references so there is no specific null Trac id value
// The value chosen here is therefore just one that will not occur in reality and is also simultaneously different from the Gitea one
// so we are more likely to detect mis-assignments.
const NullID = int64(-1)

// Accessor is the interface through which we access all Trac data.
type Accessor interface {
	/*
	 * Components
	 */
	// GetComponents retrieves all Trac components, passing each one to the provided "handler" function.
	GetComponents(handlerFn func(component *Label) error) error

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
	// GetPriorities retrieves all priorities used in Trac tickets, passing each one to the provided "handler" function.
	GetPriorities(handlerFn func(priority *Label) error) error

	/*
	 * Resolutions
	 */
	// GetResolutions retrieves all resolutions used in Trac tickets, passing each one to the provided "handler" function.
	GetResolutions(handlerFn func(resolution *Label) error) error

	/*
	 * Severities
	 */
	// GetSeverities retrieves all severities used in Trac tickets, passing each one to the provided "handler" function.
	GetSeverities(handlerFn func(severity *Label) error) error

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
	// GetTypes retrieves all types used in Trac tickets, passing each one to the provided "handler" function.
	GetTypes(handlerFn func(tracType *Label) error) error

	/*
	 * Users
	 */
	// GetUserNames retrieves the names of all users mentioned in Trac tickets, wiki pages etc., passing each one to the provided "handler" function.
	GetUserNames(handlerFn func(userName string) error) error

	/*
	 * Versions
	 */
	// GetVersions retrieves all versions used in Trac, passing each one to the provided "handler" function.
	GetVersions(handlerFn func(version *Label) error) error

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
