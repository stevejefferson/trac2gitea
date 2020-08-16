package markdown

import "regexp"

var definitionListRegexp = regexp.MustCompile(`(?m)^ ([^:]+)\:\:`)

func (converter *DefaultConverter) convertDefinitionLists(in string) string {
	return definitionListRegexp.ReplaceAllString(in, "*$1*  \n")
}
