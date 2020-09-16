// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

/*
 * Set up for ticket/issue change parts of ticket tests.
 * Contains:
 * - ticket change data types
 * - expectations for use with ticket changes.
 */

// TicketChangeImport holds the data on a ticket change import operation
type TicketChangeImport struct {
	tracChangeType trac.TicketChangeType
	issueCommentID int64
	author         *TicketUserImport
	prevOwner      *TicketUserImport
	owner          *TicketUserImport
	prevMilestone  *TicketMilestoneImport
	milestone      *TicketMilestoneImport
	prevLabel      *TicketLabelImport
	label          *TicketLabelImport
	isClose        bool
	prevSummary    string
	summary        string
	text           string
	markdownText   string
	time           int64
}

func tracTicketChangeLabelName(label *TicketLabelImport) string {
	if label == nil {
		return ""
	}

	return label.tracName
}

func createTracTicketChange(ticket *TicketImport, ticketChange *TicketChangeImport) *trac.TicketChange {
	oldValue := ""
	newValue := ""
	switch ticketChange.tracChangeType {
	case trac.TicketCommentChange:
		newValue = ticketChange.text
	case trac.TicketComponentChange:
		fallthrough
	case trac.TicketPriorityChange:
		fallthrough
	case trac.TicketResolutionChange:
		fallthrough
	case trac.TicketSeverityChange:
		fallthrough
	case trac.TicketTypeChange:
		fallthrough
	case trac.TicketVersionChange:
		oldValue = tracTicketChangeLabelName(ticketChange.prevLabel)
		newValue = tracTicketChangeLabelName(ticketChange.label)
	case trac.TicketMilestoneChange:
		oldValue = ticketChange.prevMilestone.milestoneName
		newValue = ticketChange.milestone.milestoneName
	case trac.TicketOwnerChange:
		oldValue = ticketChange.prevOwner.tracUser
		newValue = ticketChange.owner.tracUser
	case trac.TicketStatusChange:
		if ticketChange.isClose {
			newValue = trac.TicketStatusClosed
		} else {
			newValue = trac.TicketStatusReopened
		}
	case trac.TicketSummaryChange:
		oldValue = ticketChange.prevSummary
		newValue = ticketChange.summary
	}
	tracChange := trac.TicketChange{
		TicketID:   ticket.ticketID,
		ChangeType: ticketChange.tracChangeType,
		Author:     ticketChange.author.tracUser,
		OldValue:   oldValue,
		NewValue:   newValue,
		Time:       ticketChange.time,
	}

	return &tracChange
}

func expectTracChangeRetrievals(t *testing.T, ticket *TicketImport, ticketChanges ...*TicketChangeImport) {
	// expect trac accessor to return each of our trac ticket comments
	mockTracAccessor.
		EXPECT().
		GetTicketChanges(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, handlerFn func(change *trac.TicketChange) error) error {
			for _, ticketChange := range ticketChanges {
				tracChange := createTracTicketChange(ticket, ticketChange)
				handlerFn(tracChange)
			}
			return nil
		})
}
