package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	removeNonAlphaNumRegex    = regexp.MustCompile("[^ a-zA-Z0-9]")
	removeMultipleSpacesRegex = regexp.MustCompile(`\s+`)
)

func FormatUsername(name string) string {
	// remove any non-alphanumeric characters from the string
	// example "?-_!" should be ""
	// example "a?!;';';'b" should be "ab"
	validChars := removeNonAlphaNumRegex.ReplaceAllString(name, "")
	formatted := validChars

	// trim spaces
	formatted = strings.TrimSpace(formatted)

	// trim spaces between chars to maxed only one space
	// example "a       b" should be "a b"
	singleSpace := removeMultipleSpacesRegex.ReplaceAllString(formatted, " ")
	formatted = singleSpace

	// format name to lowercase
	formatted = strings.ToLower(formatted)

	// format name to replace all spaces into _ (underscore)
	formatted = strings.ReplaceAll(formatted, " ", "_")

	// generate random string id using nanoid package
	id := GenerateRandomString(7)

	formatted = fmt.Sprintf("%s_%s", formatted, id)

	return formatted
}
