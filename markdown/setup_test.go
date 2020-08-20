// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"stevejefferson.co.uk/trac2gitea/accessor/mock_gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/mock_trac"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

const (
	// random bits of text to surround Trac markdown to be converted
	// - these are used to validate that the surround context is left intact
	leadingText  = "some text"
	trailingText = "some other text"

	// name of wiki page
	wikiPage = "SomeWikiPage"
)

var ctrl *gomock.Controller
var converter *markdown.DefaultConverter
var mockTracAccessor *mock_trac.MockAccessor
var mockGiteaAccessor *mock_gitea.MockAccessor

func setUp(t *testing.T) {
	ctrl = gomock.NewController(t)

	// create mock accessors
	mockTracAccessor = mock_trac.NewMockAccessor(ctrl)
	mockGiteaAccessor = mock_gitea.NewMockAccessor(ctrl)

	// create converter to be tested
	converter = markdown.CreateWikiDefaultConverter(mockTracAccessor, mockGiteaAccessor, wikiPage)
}

func tearDown(t *testing.T) {
	ctrl.Finish()
}

func assertEquals(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("Expecting \"%v\", got \"%v\"\n", expected, got)
	}
}
