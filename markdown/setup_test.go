package markdown_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"stevejefferson.co.uk/trac2gitea/accessor/mock_gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/mock_giteawiki"
	"stevejefferson.co.uk/trac2gitea/accessor/mock_trac"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

const (
	// random bits of text to surround Trac markdown to be converted
	// - these are used to validate that the surround context is left intact
	leadingText  = "some text"
	trailingText = "some other text"
)

var ctrl *gomock.Controller
var converter *markdown.DefaultConverter
var mockTracAccessor *mock_trac.MockAccessor
var mockGiteaAccessor *mock_gitea.MockAccessor
var mockGiteaWikiAccessor *mock_giteawiki.MockAccessor

func setUp(t *testing.T) {
	ctrl = gomock.NewController(t)

	// create mock accessors
	mockTracAccessor = mock_trac.NewMockAccessor(ctrl)
	mockGiteaAccessor = mock_gitea.NewMockAccessor(ctrl)
	mockGiteaWikiAccessor = mock_giteawiki.NewMockAccessor(ctrl)

	// create converter to be tested
	converter = markdown.CreateWikiDefaultConverter(mockTracAccessor, mockGiteaAccessor, mockGiteaWikiAccessor)

}

func tearDown(t *testing.T) {
	ctrl.Finish()
}

func assertEquals(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("Expecting \"%v\", got \"%v\"\n", expected, got)
	}
}
