package internal

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToTitleCase(s string) string {
	t := cases.Title(language.English)
	return t.String(s)
}
