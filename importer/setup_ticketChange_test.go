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
	text           string
	owner          *TicketUserImport
	prevOwner      *TicketUserImport
	markdownText   string
	time           int64
}

func createTracTicketChange(ticket *TicketImport, ticketChange *TicketChangeImport) *trac.TicketChange {
	var comment *trac.TicketComment = nil
	var ownership *trac.TicketOwnership = nil
	var status *trac.TicketStatus = nil

	switch ticketChange.tracChangeType {
	case trac.TicketCommentChange:
		comment = &trac.TicketComment{Text: ticketChange.text}
	case trac.TicketOwnershipChange:
		ownership = &trac.TicketOwnership{Owner: ticketChange.owner.tracUser, PrevOwner: ticketChange.prevOwner.tracUser}
	case trac.TicketStatusChange:
		status = &trac.TicketStatus{IsClosed: true}
	}
	tracChange := trac.TicketChange{
		TicketID:   ticket.ticketID,
		Author:     ticketChange.author.tracUser,
		Time:       ticketChange.time,
		ChangeType: ticketChange.tracChangeType,
		Comment:    comment,
		Ownership:  ownership,
		Status:     status,
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
