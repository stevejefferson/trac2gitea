package markdown

import "regexp"

var definitionListRegexp = regexp.MustCompile(` ([[:alpha:]]+)\:\:\n`)

func (converter *DefaultConverter) convertDefinitionLists(in string) string {
	return definitionListRegexp.ReplaceAllString(in, "*$1*  \n")
}
