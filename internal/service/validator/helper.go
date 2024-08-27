package validator

import "regexp"

// ValidateEmail checks if the provided email address is valid
func ValidateEmail(email string) bool {
	return regexp.
		MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).
		MatchString(email)
}
