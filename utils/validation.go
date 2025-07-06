package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Error codes for consistent error handling
const (
	ErrInvalidInput        = "INVALID_INPUT"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrNotFound            = "NOT_FOUND"
	ErrDuplicateEmail      = "DUPLICATE_EMAIL"
	ErrSlotFull            = "SLOT_FULL"
	ErrDoubleBooking       = "DOUBLE_BOOKING"
	ErrInvalidDate         = "INVALID_DATE"
	ErrInvalidAge          = "INVALID_AGE"
	ErrInvalidCapacity     = "INVALID_CAPACITY"
	ErrInvalidRole         = "INVALID_ROLE"
	ErrInvalidStatus       = "INVALID_STATUS"
	ErrChildNotFound       = "CHILD_NOT_FOUND"
	ErrSlotNotFound        = "SLOT_NOT_FOUND"
	ErrReservationNotFound = "RESERVATION_NOT_FOUND"
)

// ValidationError represents a validation error with code and details
type ValidationError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(code, message string, details map[string]interface{}) ValidationError {
	return ValidationError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func IsEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsPasswordStrong(password string) bool {
	return len(password) >= 8
}

func IsFieldPresent(value string) bool {
	return len(value) > 0
}

// ValidateEmail validates email format and returns detailed error
func ValidateEmail(email string) error {
	if !IsFieldPresent(email) {
		return NewValidationError(ErrInvalidInput, "Email is required", map[string]interface{}{
			"field": "email",
			"value": email,
		})
	}
	if !IsEmailValid(email) {
		return NewValidationError(ErrInvalidInput, "Invalid email format", map[string]interface{}{
			"field": "email",
			"value": email,
		})
	}
	return nil
}

// ValidatePassword validates password strength and returns detailed error
func ValidatePassword(password string) error {
	if !IsFieldPresent(password) {
		return NewValidationError(ErrInvalidInput, "Password is required", map[string]interface{}{
			"field": "password",
		})
	}
	if !IsPasswordStrong(password) {
		return NewValidationError(ErrInvalidInput, "Password must be at least 8 characters", map[string]interface{}{
			"field":      "password",
			"min_length": 8,
		})
	}
	return nil
}

// ValidateBirthdate validates birthdate string in YYYY-MM-DD and ensures age between 0-18
func ValidateBirthdate(birthdate string) (int, error) {
	if !IsFieldPresent(birthdate) {
		return 0, NewValidationError(ErrInvalidInput, "Birthdate is required", map[string]interface{}{
			"field": "birthdate",
		})
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(birthdate) {
		return 0, NewValidationError(ErrInvalidDate, "Invalid birthdate format. Use YYYY-MM-DD", map[string]interface{}{
			"field":  "birthdate",
			"value":  birthdate,
			"format": "YYYY-MM-DD",
		})
	}
	parsed, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return 0, NewValidationError(ErrInvalidDate, "Invalid birthdate", map[string]interface{}{
			"field": "birthdate",
			"value": birthdate,
		})
	}
	today := time.Now()
	age := today.Year() - parsed.Year()
	// adjust if birthday hasn't occurred yet this year
	if today.YearDay() < parsed.YearDay() {
		age--
	}
	if age < 0 || age > 18 {
		return 0, NewValidationError(ErrInvalidAge, "Child age must be between 0 and 18", map[string]interface{}{
			"field": "birthdate",
			"value": birthdate,
		})
	}
	return age, nil
}

// ValidateChild validates child data and returns detailed error
func ValidateChild(name string, birthdate string, agePtr *int) (int, error) {
	if !IsFieldPresent(name) {
		return 0, NewValidationError(ErrInvalidInput, "Child name is required", map[string]interface{}{
			"field": "name",
		})
	}
	if len(name) > 100 {
		return 0, NewValidationError(ErrInvalidInput, "Child name is too long (max 100 characters)", map[string]interface{}{
			"field":      "name",
			"max_length": 100,
		})
	}
	if IsFieldPresent(birthdate) {
		age, err := ValidateBirthdate(birthdate)
		return age, err
	}
	if agePtr != nil {
		age := *agePtr
		if age < 0 || age > 18 {
			return 0, NewValidationError(ErrInvalidAge, "Child age must be between 0 and 18", map[string]interface{}{
				"field": "age",
				"value": age,
			})
		}
		return age, nil
	}
	return 0, NewValidationError(ErrInvalidInput, "Either birthdate or age is required", nil)
}

// ValidateSlot validates slot data and returns detailed error
func ValidateSlot(date string, capacity int) error {
	if !IsFieldPresent(date) {
		return NewValidationError(ErrInvalidInput, "Slot date is required", map[string]interface{}{
			"field": "date",
			"value": date,
		})
	}

	// Validate date format (YYYY-MM-DD)
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(date) {
		return NewValidationError(ErrInvalidDate, "Invalid date format. Use YYYY-MM-DD", map[string]interface{}{
			"field":  "date",
			"value":  date,
			"format": "YYYY-MM-DD",
		})
	}

	// Validate date is not in the past
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return NewValidationError(ErrInvalidDate, "Invalid date", map[string]interface{}{
			"field": "date",
			"value": date,
		})
	}

	if parsedDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return NewValidationError(ErrInvalidDate, "Slot date cannot be in the past", map[string]interface{}{
			"field": "date",
			"value": date,
		})
	}

	if capacity <= 0 || capacity > 100 {
		return NewValidationError(ErrInvalidCapacity, "Slot capacity must be between 1 and 100", map[string]interface{}{
			"field": "capacity",
			"value": capacity,
			"min":   1,
			"max":   100,
		})
	}

	return nil
}

// ValidateReservation validates reservation data and returns detailed error
func ValidateReservation(childID, slotID uint) error {
	if childID == 0 {
		return NewValidationError(ErrInvalidInput, "Child ID is required", map[string]interface{}{
			"field": "child_id",
			"value": childID,
		})
	}
	if slotID == 0 {
		return NewValidationError(ErrInvalidInput, "Slot ID is required", map[string]interface{}{
			"field": "slot_id",
			"value": slotID,
		})
	}
	return nil
}

// ValidateRole validates user role and returns detailed error
func ValidateRole(role string) error {
	validRoles := []string{"parent", "admin"}
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}
	return NewValidationError(ErrInvalidRole, "Invalid role. Must be 'parent' or 'admin'", map[string]interface{}{
		"field":       "role",
		"value":       role,
		"valid_roles": validRoles,
	})
}

// ValidateStatus validates reservation status and returns detailed error
func ValidateStatus(status string) error {
	validStatuses := []string{"pending", "approved", "rejected"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return NewValidationError(ErrInvalidStatus, "Invalid status. Must be 'pending', 'approved', or 'rejected'", map[string]interface{}{
		"field":          "status",
		"value":          status,
		"valid_statuses": validStatuses,
	})
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, perPage int) error {
	if page < 1 {
		return NewValidationError(ErrInvalidInput, "Page must be greater than 0", map[string]interface{}{
			"field": "page",
			"value": page,
			"min":   1,
		})
	}
	if perPage < 1 || perPage > 100 {
		return NewValidationError(ErrInvalidInput, "Per page must be between 1 and 100", map[string]interface{}{
			"field": "per_page",
			"value": perPage,
			"min":   1,
			"max":   100,
		})
	}
	return nil
}

// ParsePagination parses and validates pagination parameters from query string
func ParsePagination(pageStr, perPageStr string) (page, perPage int, err error) {
	page = 1
	perPage = 20

	if pageStr != "" {
		if page, err = strconv.Atoi(pageStr); err != nil {
			return 0, 0, NewValidationError(ErrInvalidInput, "Invalid page number", map[string]interface{}{
				"field": "page",
				"value": pageStr,
			})
		}
	}

	if perPageStr != "" {
		if perPage, err = strconv.Atoi(perPageStr); err != nil {
			return 0, 0, NewValidationError(ErrInvalidInput, "Invalid per_page number", map[string]interface{}{
				"field": "per_page",
				"value": perPageStr,
			})
		}
	}

	if err := ValidatePagination(page, perPage); err != nil {
		return 0, 0, err
	}

	return page, perPage, nil
}

// ParseUint parses a string to uint
func ParseUint(s string) (uint, error) {
	var result uint
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
