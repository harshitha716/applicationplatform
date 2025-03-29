package helper

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// IsValidEmail checks if the provided email is valid
func IsValidEmail(email string) bool {
	// Define a regular expression for validating an email address
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	// Match the email with the regex pattern
	return emailRegex.MatchString(SanitizeEmail(email))
}

func IsZampEmail(email string) bool {
	// Define a regular expression for validating an email address
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@zamp.ai$`)

	// Match the email with the regex pattern
	return emailRegex.MatchString(email)
}

func GetDomainFromEmail(email string) string {
	// Extract the domain from the email
	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func AreEmailsEqual(email1 string, email2 string) bool {
	return SanitizeEmail(email1) == SanitizeEmail(email2)
}

func GetNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) < 1 {
		return ""
	}

	// split on "." and take the first part and capitalize it
	parts = strings.Split(parts[0], ".")
	if len(parts) < 1 {
		return ""
	}
	return cases.Title(language.Und).String(parts[0])
}
