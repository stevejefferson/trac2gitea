package markdown_test

import (
	"testing"

	"stevejefferson.co.uk/trac2gitea/markdown"
)

var converter *markdown.DefaultConverter

// mockGiteaWikiAccessor is a mock implementation of accessor.giteawiki.Accessor
type mockGiteaWikiAccessor struct {
}

// mockTracAccessor is a mock implementation of accessor.trac.Accessor
type mockTracAccessor struct {
}

func setUp() {
	// tracAccessor := mockTracAccessor{}
	// giteaAccessor := mockGiteaAccessor{}
	// giteaWikiAccessor := mockGiteaWikiAccessor{}

	// converter = markdown.CreateWikiDefaultConverter(tracAccessor, giteaAccessor, giteaWikiAccessor)
}

func TestAnchor(t *testing.T) {
}
