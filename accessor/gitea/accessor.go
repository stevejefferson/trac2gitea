// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import "database/sql"

// Accessor is the interface to all of our interactions with a Gitea project.
type Accessor interface {
	/*
	 * Attachments
	 */
	// GetAttachmentUUID returns the UUID for a named attachment of a given issue - returns empty string if cannot find issue/attachment.
	GetAttachmentUUID(issueID int64, name string) (string, error)

	// AddAttachment adds a new attachment to a given issue with the provided data - returns id of created attachment
	AddAttachment(uuid string, issueID int64, commentID int64, attachmentName string, attachmentFile string, time int64) (int64, error)

	// GetAttachmentURL retrieves the URL for viewing a Gitea attachment
	GetAttachmentURL(uuid string) string

	/*
	 * Comments
	 */
	// AddComment adds a comment to Gitea - returns id of created comment.
	AddComment(issueID int64, authorID int64, comment string, time int64) (int64, error)

	// GetCommentURL retrieves the URL for viewing a Gitea comment for a given issue.
	GetCommentURL(issueID int64, commentID int64) string

	// GetCommentID retrives the ID of a given comment for a given issue or -1 if no such issue/comment
	GetCommentID(issueID int64, commentStr string) (int64, error)

	/*
	 * Configuration
	 */
	// GetStringConfig retrieves a value from the Gitea config as a string.
	GetStringConfig(sectionName string, configName string) string

	/*
	 * Issues
	 */
	// GetIssueID retrieves the id of the Gitea issue corresponding to a given Trac ticket - returns -1 if no such issue.
	GetIssueID(ticketID int64) (int64, error)

	// AddIssue adds a new issue to Gitea - returns id of created issue.
	AddIssue(
		ticketID int64,
		summary string,
		reporterID int64,
		milestone string,
		ownerID sql.NullString,
		owner string,
		closed bool,
		description string,
		created int64) (int64, error)

	// SetIssueUpdateTime sets the update time on a given Gitea issue.
	SetIssueUpdateTime(issueID int64, updateTime int64) error

	// GetIssueURL retrieves a URL for viewing a given issue
	GetIssueURL(issueID int64) string

	/*
	 * Issue Labels
	 */
	// GetIssueLabelID retrieves the id of the given Gitea issue and label - returns -1 if no such issue label.
	GetIssueLabelID(issueID int64, labelID int64) (int64, error)

	// AddIssueLabel adds an issue label to Gitea, returns issue label ID
	AddIssueLabel(issueID int64, labelID int64) (int64, error)

	/*
	 * Labels
	 */
	// GetLabelID retrieves the id of the given label, returns -1 if no such label
	GetLabelID(labelName string) (int64, error)

	// AddLabel adds a label to Gitea, returns label id.
	AddLabel(label string, color string) (int64, error)

	/*
	 * Milestones
	 */
	// AddMilestone adds a milestone to Gitea,  returns id of created milestone
	AddMilestone(name string, content string, closed bool, deadlineTimestamp int64, closedTimestamp int64) (int64, error)

	// GetMilestoneID gets the ID of a named milestone - returns -1 if no such milestone
	GetMilestoneID(name string) (int64, error)

	// GetMilestoneURL gets the URL for accessing a given milestone
	GetMilestoneURL(milestoneID int64) string

	/*
	 * Repository
	 */
	// UpdateRepoIssueCount updates the count of total and closed issue for a our chosen Gitea repository.
	UpdateRepoIssueCount(count int, closedCount int) error

	// GetCommitURL retrieves the URL for viewing a given commit in the current repository
	GetCommitURL(commitID string) string

	// GetSourceURL retrieves the URL for viewing the latest version of a source file on a given branch of the current repository
	GetSourceURL(branchPath string, filePath string) string

	/*
	 * Users
	 */
	// GetUserID retrieves the id of a named Gitea user - returns -1 if no such user.
	GetUserID(name string) (int64, error)

	// GetDefaultAssigneeID retrieves the id of the user to which to assign tickets/comments in the case where the Trac-supplied user id does not exist in Gitea.
	GetDefaultAssigneeID() int64

	// GetDefaultAuthorID retrieves the id of the user to set as the author of tickets/comments in the case where the Trac-supplied user id does not exist in Gitea.
	GetDefaultAuthorID() int64

	// GetUserEMailAddress retrieves the email address of a given user
	GetUserEMailAddress(userID int64) (string, error)

	/*
	 * Wiki
	 */
	// GetWikiAttachmentRelPath returns the location of an attachment to Trac a wiki page when stored in the Gitea wiki repository.
	// The returned path is relative to the root of the Gitea wiki repository.
	GetWikiAttachmentRelPath(pageName string, filename string) string

	// GetWikiHtdocRelPath returns the location of a given Trac 'htdocs' file when stored in the Gitea wiki repository.
	// The returned path is relative to the root of the Gitea wiki repository.
	GetWikiHtdocRelPath(filename string) string

	// GetWikiFileURL returns a URL for viewing a file stored in the Gitea wiki repository.
	GetWikiFileURL(relpath string) string

	// CloneWiki clones the wiki repo.
	CloneWiki() error

	// LogWiki returns the log of commits for the given wiki page.
	LogWiki(pageName string) ([]string, error)

	// CommitWiki commits any files added or updated since the last commit to our cloned wiki repo.
	CommitWiki(author string, authorEMail string, message string) error

	// PushWiki pushes all changes to the local wiki repository back to the remote.
	PushWiki() error

	// CopyFileToWiki copies an external file into the local clone of the Gitea Wiki
	CopyFileToWiki(externalFilePath string, giteaWikiRelPath string) error

	// WriteWikiPage writes (a version of) a wiki page to the local clone of the wiki repository, returning the path to the written file.
	WriteWikiPage(pageName string, markdownText string) (string, error)

	// TranslateWikiPageName translates a Trac wiki page name into a Gitea one
	TranslateWikiPageName(pageName string) string
}
