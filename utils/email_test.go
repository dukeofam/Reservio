package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMail_Mock(t *testing.T) {
	// This is a placeholder test to show SendMail can be called.
	// In real tests, use a mock or interface to avoid sending real emails.
	err := SendMail("to@example.com", "Test Subject", "Test Body")
	// We expect an error because credentials are fake, but function should be callable
	assert.Error(t, err)
}
