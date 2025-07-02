package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	assert.True(t, IsEmailValid("test@example.com"))
	assert.False(t, IsEmailValid("notanemail"))
}

func TestIsPasswordStrong(t *testing.T) {
	assert.True(t, IsPasswordStrong("12345678"))
	assert.False(t, IsPasswordStrong("short"))
}

func TestIsFieldPresent(t *testing.T) {
	assert.True(t, IsFieldPresent("hello"))
	assert.False(t, IsFieldPresent(""))
}
