package trac

// Accessor is the interface through which we access all Trac data.
type Accessor interface {
	/*
	 * Attachments
	 */
	// GetAttachmentPath retrieves the path to a named attachment to a Trac ticket.
	GetAttachmentPath(ticketID int64, name string) string

	// GetAttachments retrieves all attachments for a given Trac ticket, passing data from each one to the provided "handler" function.
	GetAttachments(ticketID int64, handlerFn func(ticketID int64, time int64, size int64, author string, filename string, description string))

	/*
	 * Comments
	 */
	// GetComments retrieves all comments on a given Trac ticket, passing data from each one to the provided "handler" function.
	GetComments(ticketID int64, handlerFn func(ticketID int64, time int64, author string, comment string))

	// GetCommentString retrieves a given comment string for a given Trac ticket
	GetCommentString(ticketID int64, commentNum int64) string

	/*
	 * Components
	 */
	// GetComponentNames retrieves all Trac component names, passing each one to the provided "handler" function.
	GetComponentNames(handlerFn func(cmptName string))

	/*
	 * Configuration
	 */
	// GetStringConfig retrieves a value from the Trac config as a string.
	GetStringConfig(sectionName string, configName string) string

	/*
	 * Milestones
	 */
	// GetMilestones retrieves all Trac milestones, passing data from each one to the provided "handler" function.
	GetMilestones(handlerFn func(name string, description string, due int64, completed int64))

	/*
	 * Paths
	 */
	// GetFullPath retrieves the absolute path of a path relative to the root of the Trac installation.
	GetFullPath(element ...string) string

	/*
	 * Priorities
	 */
	// GetPriorityNames retrieves all priority names used in Trac tickets, passing each one to the provided "handler" function.
	GetPriorityNames(handlerFn func(priorityName string))

	/*
	 * Resolutions
	 */
	// GetResolutionNames retrieves all resolution names used in Trac tickets, passing each one to the provided "handler" function.
	GetResolutionNames(handlerFn func(resolution string))

	/*
	 * Severities
	 */
	// GetSeverityNames retrieves all severity names used in Trac tickets, passing each one to the provided "handler" function.
	GetSeverityNames(handlerFn func(severityName string))

	/*
	 * Tickets
	 */
	// GetTickets retrieves all Trac tickets, passing data from each one to the provided "handler" function.
	GetTickets(handlerFn func(
		ticketID int64, ticketType string, created int64,
		component string, severity string, priority string,
		owner string, reporter string, version string,
		milestone string, status string, resolution string,
		summary string, description string))

	/*
	 * Types
	 */
	// GetTypeNames retrieves all type names used in Trac tickets, passing each one to the provided "handler" function.
	GetTypeNames(handlerFn func(typeName string))

	/*
	 * Versions
	 */
	// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
	GetVersionNames(handlerFn func(version string))

	/*
	 * Wiki
	 */
	// GetWikiPages retrieves all Trac wiki pages, passing data from each one to the provided "handler" function.
	GetWikiPages(handlerFn func(pageName string, pageText string, author string, comment string, version int64, updateTime int64))

	// IsPredefinedPage returns true if the provided page name is one of Trac's predefined ones - by default we ignore these
	IsPredefinedPage(pageName string) bool
}
