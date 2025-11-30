package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type (
	// Validator - validation helper.
	Validator struct {
		// FieldErrors - Form fields errors.
		FieldErrors map[string]string
		// NonFieldErrors - error not related to form fields.
		NonFieldErrors []string
	}
	// ValidationFunction - validation function.
	ValidationFunction[T comparable] func(T) bool
)

// EmailRX - email validation regular expression.
var EmailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+" +
		"@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
		"(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

// Valid - check if validation succeed.
func (validator *Validator) Valid() bool {
	return len(validator.FieldErrors) == 0 && len(validator.NonFieldErrors) == 0
}

// AddFieldError - add error for field after validation.
func (validator *Validator) AddFieldError(key, message string) {
	if validator.FieldErrors == nil {
		validator.FieldErrors = make(map[string]string)
	}

	_, exists := validator.FieldErrors[key]

	if !exists {
		validator.FieldErrors[key] = message
	}
}

// AddNonFieldError - add errors not related to form fields.
func (validator *Validator) AddNonFieldError(message string) {
	validator.NonFieldErrors = append(validator.NonFieldErrors, message)
}

// CheckField - checks if field is valid using validator and
// validation function.
func CheckField[T comparable](
	validator *Validator,
	validationFunc ValidationFunction[T],
	value T, key, message string,
) {
	valid := validationFunc(value)

	if !valid {
		validator.AddFieldError(key, message)
	}
}

// CreateNotBlankValidator - not blank field validation.
func CreateNotBlankValidator() ValidationFunction[string] {
	return func(value string) bool {
		return strings.TrimSpace(value) != ""
	}
}

// CreateMaxCharsValidator - maximum characters in field validation.
func CreateMaxCharsValidator(limit int) ValidationFunction[string] {
	return func(value string) bool {
		return utf8.RuneCountInString(value) <= limit
	}
}

// CreatePermittedValueValidator - checks that field value is within provided
// permitted values.
func CreatePermittedValueValidator[T comparable](permittedValues ...T) ValidationFunction[T] {
	return func(value T) bool {
		return slices.Contains(permittedValues, value)
	}
}

// CreateMinCharsValidator - minimum characters in field validation.
func CreateMinCharsValidator(limit int) ValidationFunction[string] {
	return func(value string) bool {
		return utf8.RuneCountInString(value) >= limit
	}
}

// CreateMatchesRegexValidator - checks that provided string value matches provided
// regular expression.
func CreateMatchesRegexValidator(regex *regexp.Regexp) ValidationFunction[string] {
	return func(value string) bool {
		return regex.MatchString(value)
	}
}
