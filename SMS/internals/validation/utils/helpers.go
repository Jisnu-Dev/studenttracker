package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var NameRegex = regexp.MustCompile(`^[\p{L}]+(?: [\p{L}]+)*$`)

var AlphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

const (
	MaxEmailLength    = 254
	MaxPasswordLength = 72
)

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func ContainsTab(s string) bool {
	return strings.ContainsRune(s, '\t')
}

func HasLeadingOrTrailingSpace(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(s)
	return unicode.IsSpace(runes[0]) || unicode.IsSpace(runes[len(runes)-1])
}

func HasMultipleConsecutiveSpaces(s string) bool {
	return strings.Contains(s, "  ")
}

func ValidateCleanWhitespace(fieldName, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if IsBlank(value) {
		return fmt.Errorf("%s cannot be empty or only whitespace", fieldName)
	}
	if ContainsTab(value) {
		return fmt.Errorf("%s cannot contain tab characters", fieldName)
	}
	if HasLeadingOrTrailingSpace(value) {
		return fmt.Errorf("%s cannot have leading or trailing spaces", fieldName)
	}
	if HasMultipleConsecutiveSpaces(value) {
		return fmt.Errorf("%s cannot contain multiple consecutive spaces", fieldName)
	}
	return nil
}

type ValidationErrors map[string][]string

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}
	parts := make([]string, 0, len(v))
	for field, msgs := range v {
		for _, msg := range msgs {
			parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
		}
	}
	return strings.Join(parts, "; ")
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func (v ValidationErrors) Add(field string, err error) {
	if err == nil {
		return
	}
	v[field] = append(v[field], err.Error())
}

func (v ValidationErrors) AsError() error {
	if v.HasErrors() {
		return v
	}
	return nil
}
