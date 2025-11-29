package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateSnippetCreateForm(form *snippetCreateForm) {
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = ValidationErrorBlank
	} else if utf8.RuneCountInString(form.Title) > TitleLengthLimit {
		form.FieldErrors["title"] = fmt.Sprintf("This field cannot be more than %d characters long", TitleLengthLimit)
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = ValidationErrorBlank
	}

	if form.Expires != ExpiresInDay && form.Expires != ExpiresInWeek && form.Expires != ExpiresInYear {
		form.FieldErrors["expires"] = fmt.Sprintf(
			"This field must be either %d, %d or %d",
			ExpiresInDay,
			ExpiresInWeek,
			ExpiresInYear,
		)
	}
}
