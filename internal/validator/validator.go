package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type (
	Validator struct {
		FieldErrors map[string]string
	}
	ValidationFunction[T comparable] func(T) bool
)

func (validator *Validator) Valid() bool {
	return len(validator.FieldErrors) == 0
}

func (validator *Validator) AddFieldError(key, message string) {
	if validator.FieldErrors == nil {
		validator.FieldErrors = make(map[string]string)
	}

	_, exists := validator.FieldErrors[key]

	if !exists {
		validator.FieldErrors[key] = message
	}
}

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

func CreateNotBlankValidator() ValidationFunction[string] {
	return func(value string) bool {
		return strings.TrimSpace(value) != ""
	}
}

func CreateMaxCharsValidator(limit int) ValidationFunction[string] {
	return func(value string) bool {
		return utf8.RuneCountInString(value) <= limit
	}
}

func CreatePermittedValueValidator[T comparable](permittedValues ...T) ValidationFunction[T] {
	return func(value T) bool {
		return slices.Contains(permittedValues, value)
	}
}
