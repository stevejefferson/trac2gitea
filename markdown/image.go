package markdown

import (
	"fmt"
	"regexp"
)

func (converter *Converter) convertImageReference(in string, linkprefix string) string {
	regex := regexp.MustCompile(`\[\[Image\(([^,\)]*)[^\)]*\)\]\]`)
	out := regex.ReplaceAllStringFunc(in, func(m string) string {
		u := regex.ReplaceAllString(m, "$1")
		u = converter.resolveTracLink(u, linkprefix)
		return fmt.Sprintf("![](%s)", u)
	})

	return out
}
