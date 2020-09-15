// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"os"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// getCommentIssueComment returns a Gitea issue comment reflecting a comment made on ta Trac ticket
func (importer *Importer) getCommentIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string) (*gitea.IssueComment, error) {
	convertedText := importer.markdownConverter.TicketConvert(change.TicketID, change.Comment.Text)
	giteaComment := gitea.IssueComment{CommentType: gitea.CommentIssueCommentType, Text: convertedText}

	return &giteaComment, nil
}

// getOwnershipIssueComment returns a Gitea issue comment reflecting an ownership change return nil if we cannot express ownership change in Gitea
func (importer *Importer) getOwnershipIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string) (*gitea.IssueComment, error) {
	var err error
	prevOwnerID := int64(0)
	prevOwnerName := change.Ownership.PrevOwner
	if prevOwnerName != "" {
		prevOwnerID, err = importer.getUser(prevOwnerName, userMap)
		if err != nil {
			return nil, err
		}
		if prevOwnerID == -1 {
			return nil, nil // cannot map user onto Gitea
		}
	}

	assigneeID := int64(0)
	removedAssigneeID := int64(0)
	ownerName := change.Ownership.Owner
	if ownerName != "" {
		assigneeID, err = importer.getUser(ownerName, userMap)
		if err != nil {
			return nil, err
		}
		if assigneeID == -1 {
			return nil, nil // cannot map user onto Gitea
		}
	} else {
		// this is an assignee removal
		removedAssigneeID = prevOwnerID
	}

	if issueID == 1 {
		fmt.Fprintf(os.Stderr, "XXX issueID=%d, prevOwnerName=%s, ownerName=%s, assigneeID=%d, removedAssigneeID=%d\n", issueID, prevOwnerName, ownerName, assigneeID, removedAssigneeID)
	}

	giteaComment := gitea.IssueComment{
		CommentType:     gitea.AssigneeIssueCommentType,
		AssigneeID:      assigneeID,
		RemovedAssignee: removedAssigneeID,
	}
	return &giteaComment, nil
}

// importTicketChange imports a single ticket change from Trac to Gitea, returns ID of created Gitea comment or -1 if comment already exists
func (importer *Importer) importTicketChange(issueID int64, change *trac.TicketChange, userMap map[string]string) (int64, error) {
	var issueComment *gitea.IssueComment
	var err error

	// obtain Gitea comment to add
	switch change.ChangeType {
	case trac.TicketCommentChange:
		issueComment, err = importer.getCommentIssueComment(issueID, change, userMap)
	case trac.TicketOwnershipChange:
		issueComment, err = importer.getOwnershipIssueComment(issueID, change, userMap)
	}
	if err != nil {
		return -1, err
	}

	// eliminate cases where we could not map the Trac ticket change onto a Gitea comment
	if issueComment == nil {
		return -1, nil
	}

	// record Trac change author as original author if it cannot be mapped onto a Gitea user
	authorID, err := importer.getUser(change.Author, userMap)
	if err != nil {
		return -1, err
	}
	originalAuthorName := ""
	if authorID == -1 {
		authorID = importer.defaultAuthorID
		originalAuthorName = change.Author
	}
	issueComment.AuthorID = authorID
	issueComment.OriginalAuthorID = 0
	issueComment.OriginalAuthorName = originalAuthorName

	issueComment.Time = change.Time

	// add Gitea issue comment
	commentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return -1, err
	}

	// add association between issue and comment author
	err = importer.giteaAccessor.AddIssueUser(issueID, authorID)
	if err != nil {
		return -1, err
	}

	return commentID, err
}

func (importer *Importer) importTicketChanges(ticketID int64, issueID int64, lastUpdate int64, userMap map[string]string) (int64, error) {
	commentLastUpdate := lastUpdate
	err := importer.tracAccessor.GetTicketChanges(ticketID, func(change *trac.TicketChange) error {
		commentID, err := importer.importTicketChange(issueID, change, userMap)
		if err != nil {
			return err
		}
		if commentID != -1 && commentLastUpdate < change.Time {
			commentLastUpdate = change.Time
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return commentLastUpdate, nil
}
