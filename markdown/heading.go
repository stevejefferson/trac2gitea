package markdown

import (
	"regexp"
)

func (converter *Converter) convertDelimitedHeadings(in string, tracDelimiter string, markdownDelimiter string) string {
	if len(tracDelimiter) == 0 {
		return in
	}

	tracDelimiterChar := tracDelimiter[0:1]

	//regexStr := `^` + tracDelimiter + `([^` + tracDelimiterChar + `]+)` + tracDelimiter + `$`
	regexStr := tracDelimiter + `([^` + tracDelimiterChar + `]+)` + tracDelimiter
	regex := regexp.MustCompile(regexStr)

	replacementStr := markdownDelimiter + `$1`
	out := regex.ReplaceAllString(in, replacementStr)

	// recurse to next level of heading
	return converter.convertDelimitedHeadings(out, tracDelimiter[1:], markdownDelimiter[1:])
}

func (converter *Converter) convertHeading(in string) string {
	// recurse through all heading levels starting from longest (doing shortest first risks premature regexp matches)
	return converter.convertDelimitedHeadings(in, "=======", "#######")
}
