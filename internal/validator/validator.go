package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type (
	// Validator - validation helper.
	Validator struct {
		FieldErrors map[string]string
	}
	// ValidationFunction - validation function.
	ValidationFunction[T comparable] func(T) bool
)

// Valid - check if validation succeed.
func (validator *Validator) Valid() bool {
	return len(validator.FieldErrors) == 0
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
