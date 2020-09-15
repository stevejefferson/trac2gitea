// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

// Issue describes a Gitea issue.
type Issue struct {
	Index              int64
	Summary            string
	ReporterID         int64
	Milestone          string
	OriginalAuthorID   int64
	OriginalAuthorName string
	Closed             bool
	Description        string
	Created            int64
	Updated            int64
}

// IssueAttachment describes an attachment to a Gitea issue.
type IssueAttachment struct {
	UUID      string
	CommentID int64
	FileName  string
	Time      int64
}

// IssueCommentType defines the types of issue comment we support
type IssueCommentType int64

const (
	// CommentIssueCommentType is an IssueComment reflecting a comment
	CommentIssueCommentType IssueCommentType = 0

	// AssigneeIssueCommentType is an IssueComment reflecting an assignee change
	AssigneeIssueCommentType IssueCommentType = 9
)

// IssueComment describes a comment on a Gitea issue.
type IssueComment struct {
	CommentType        IssueCommentType
	AuthorID           int64
	OriginalAuthorID   int64
	OriginalAuthorName string
	AssigneeID         int64
	RemovedAssignee    int64
	Text               string
	Time               int64
}

// Milestone describes a Gitea milestone.
type Milestone struct {
	Name        string
	Description string
	Closed      bool
	DueTime     int64
	ClosedTime  int64
}

// Accessor is the interface to all of our interactions with a Gitea project.
type Accessor interface {
	/*
	 * Configuration
	 */
	// GetStringConfig retrieves a value from the Gitea config as a string.
	GetStringConfig(sectionName string, configName string) string

	/*
	 * Issues
	 */
	// GetIssueID retrieves the id of the Gitea issue corresponding to a given index - returns -1 if no such issue.
	GetIssueID(issueIndex int64) (int64, error)

	// AddIssue adds a new issue to Gitea - returns id of created issue.
	AddIssue(issue *Issue) (int64, error)

	// SetIssueUpdateTime sets the update time on a given Gitea issue.
	SetIssueUpdateTime(issueID int64, updateTime int64) error

	// GetIssueURL retrieves a URL for viewing a given issue
	GetIssueURL(issueID int64) string

	// UpdateIssueCommentCount updates the count of comments a given issue
	UpdateIssueCommentCount(issueID int64) error

	/*
	 * Issue Assignees
	 */
	// AddIssueAssignee adds an assignee to a Gitea issue
	AddIssueAssignee(issueID int64, assigneeID int64) error

	/*
	 * Issue Attachments
	 */
	// GetIssueAttachmentUUID returns the UUID for a named attachment of a given issue - returns empty string if cannot find issue/attachment.
	GetIssueAttachmentUUID(issueID int64, fileName string) (string, error)

	// AddIssueAttachment adds a new attachment to an issue using the provided file - returns id of created attachment
	AddIssueAttachment(issueID int64, attachment *IssueAttachment, filePath string) (int64, error)

	// GetIssueAttachmentURL retrieves the URL for viewing a Gitea attachment
	GetIssueAttachmentURL(issueID int64, uuid string) string

	/*
	 * Issue Comments
	 */
	// GetIssueCommentIDsByTime retrieves the IDs of all comments created at a given time for a given issue
	GetIssueCommentIDsByTime(issueID int64, createdTime int64) ([]int64, error)

	// AddIssueComment adds a comment on a Gitea issue, returns id of created comment
	AddIssueComment(issueID int64, comment *IssueComment) (int64, error)

	// GetIssueCommentURL retrieves the URL for viewing a Gitea comment for a given issue.
	GetIssueCommentURL(issueID int64, commentID int64) string

	/*
	 * Issue Labels
	 */
	// AddIssueLabel adds an issue label to Gitea, returns issue label ID
	AddIssueLabel(issueID int64, labelID int64) (int64, error)

	/*
	 * Issue Users
	 */
	// AddIssueUser adds a user as being associated with a Gitea issue
	AddIssueUser(issueID int64, userID int64) error

	/*
	 * Labels
	 */
	// AddLabel adds a label to Gitea, returns label id.
	AddLabel(label string, color string) (int64, error)

	/*
	 * Milestones
	 */
	// GetMilestoneID gets the ID of a named milestone - returns -1 if no such milestone
	GetMilestoneID(name string) (int64, error)

	// AddMilestone adds a milestone to Gitea,  returns id of created milestone
	AddMilestone(milestone *Milestone) (int64, error)

	// GetMilestoneURL gets the URL for accessing a given milestone
	GetMilestoneURL(milestoneID int64) string

	/*
	 * Repository
	 */
	// UpdateRepoIssueCounts updates issue counts for a our chosen Gitea repository.
	UpdateRepoIssueCounts() error

	// UpdateRepoMilestoneCounts updates milestone counts for a our chosen Gitea repository.
	UpdateRepoMilestoneCounts() error

	// GetCommitURL retrieves the URL for viewing a given commit in the current repository
	GetCommitURL(commitID string) string

	// GetSourceURL retrieves the URL for viewing the latest version of a source file on a given branch of the current repository
	GetSourceURL(branchPath string, filePath string) string

	/*
	 * Users
	 */
	// GetUserID retrieves the id of a named Gitea user - returns -1 if no such user.
	GetUserID(userName string) (int64, error)

	// GetUserEMailAddress retrieves the email address of a given user
	GetUserEMailAddress(userName string) (string, error)

	// MatchUser retrieves the name of the user best matching a user name or email address
	MatchUser(userName string, userEmail string) (string, error)

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

	// CloneWiki creates a local clone of the wiki repo.
	CloneWiki() error

	// CommitWiki commits any files added or updated since the last commit to our local wiki repo.
	CommitWiki(author string, authorEMail string, message string) error

	// PushWiki pushes all changes to the local wiki repository back to the remote.
	PushWiki() error

	// CopyFileToWiki copies an external file into the local clone of the Gitea Wiki
	CopyFileToWiki(externalFilePath string, giteaWikiRelPath string) error

	// WriteWikiPage potentially writes a wiki page to the local wiki repository, returning a flag to say whether the file was physically written.
	// If a previous commit of the wiki page is found containing the provided marker string then the page will only be written if an explicit override has been provided.
	WriteWikiPage(pageName string, markdownText string, commitMarker string) (bool, error)

	// TranslateWikiPageName translates a Trac wiki page name into a Gitea one
	TranslateWikiPageName(pageName string) string
}
