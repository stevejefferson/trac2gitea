// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stevejefferson/trac2gitea/accessor/mock_gitea"
	"github.com/stevejefferson/trac2gitea/accessor/mock_trac"
	"github.com/stevejefferson/trac2gitea/importer"
	"github.com/stevejefferson/trac2gitea/mock_markdown"
)

var ctrl *gomock.Controller
var dataImporter *importer.Importer
var mockTracAccessor *mock_trac.MockAccessor
var mockGiteaAccessor *mock_gitea.MockAccessor
var mockMarkdownConverter *mock_markdown.MockConverter

func setUp(t *testing.T) {
	ctrl = gomock.NewController(t)

	// create mocks
	mockTracAccessor = mock_trac.NewMockAccessor(ctrl)
	mockGiteaAccessor = mock_gitea.NewMockAccessor(ctrl)
	mockMarkdownConverter = mock_markdown.NewMockConverter(ctrl)

	// create importer to be tested
	dataImporter, _ = importer.CreateImporter(mockTracAccessor, mockGiteaAccessor, mockMarkdownConverter, false)
}

func tearDown(t *testing.T) {
	ctrl.Finish()
}

func assertTrue(t *testing.T, assertion bool) {
	if !assertion {
		t.Errorf("Assertion failed!\n")
	}
}

func assertEquals(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("Expecting \"%v\", got \"%v\"\n", expected, got)
		//debug.PrintStack()
	}
}
