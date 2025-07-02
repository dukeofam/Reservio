package utils

import (
	"regexp"
)

func IsEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsPasswordStrong(password string) bool {
	return len(password) >= 8
}

func IsFieldPresent(value string) bool {
	return len(value) > 0
}

// TODO: Implement
