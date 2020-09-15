// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
)

/*
 * Set up for ticket/issue label parts of ticket tests.
 * Contains:
 * - ticket label data types
 * - expectations for use with ticket labels.
 */

// TicketLabelImport holds the data on a label associated with an imported ticket
type TicketLabelImport struct {
	name         string
	labelName    string
	labelID      int64
	issueLabelID int64
}

func createTicketLabelImport(prefix string, ticketLabelMap map[string]string) *TicketLabelImport {
	ticketLabel := TicketLabelImport{
		name:         prefix + "-name",
		labelName:    prefix + "-label",
		labelID:      allocateID(),
		issueLabelID: allocateID(),
	}

	ticketLabelMap[ticketLabel.name] = ticketLabel.labelName
	return &ticketLabel
}

func expectLabelRetrieval(t *testing.T, label *TicketLabelImport) {
	mockGiteaAccessor.
		EXPECT().
		AddLabel(gomock.Eq(label.labelName), gomock.Any()).
		Return(label.labelID, nil)
}

func expectIssueLabelRetrieval(t *testing.T, ticket *TicketImport, ticketLabel *TicketLabelImport) {
	// expect retrieval/creation of underlying label first
	expectLabelRetrieval(t, ticketLabel)

	mockGiteaAccessor.
		EXPECT().
		AddIssueLabel(gomock.Eq(ticket.issueID), gomock.Eq(ticketLabel.labelID)).
		Return(ticketLabel.issueLabelID, nil)
}
