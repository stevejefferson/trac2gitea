package markdown

// Converter is the interface for Trac markdown to Gitea markdown conversions
type Converter interface {
	// Convert converts a string of Trac markdown to Gitea markdown
	Convert(in string) string
}
