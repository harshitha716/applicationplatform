package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	assert.True(t, IsValidEmail("test@example.com"))
	assert.True(t, IsValidEmail("Test@example.com"))
	assert.True(t, IsValidEmail("TEST@example.com "))
	assert.False(t, IsValidEmail("invalid-email"))
}

func TestSanitizeEmail(t *testing.T) {
	assert.Equal(t, "test@example.com", SanitizeEmail(" Test@example.com "))
	assert.Equal(t, "test@example.com", SanitizeEmail("test@example.com"))
	assert.Equal(t, "test@example.com", SanitizeEmail("test@example.com"))
}

func TestAreEmailsEqual(t *testing.T) {
	assert.True(t, AreEmailsEqual("test@example.com", "test@example.com"))
	assert.True(t, AreEmailsEqual("test@example.com", "Test@example.com"))
	assert.True(t, AreEmailsEqual("test@example.com", "TEST@example.com "))
	assert.False(t, AreEmailsEqual("test@example.com", "test@example.com.com"))
}

func TestGetDomainFromEmail(t *testing.T) {
	assert.Equal(t, "example.com", GetDomainFromEmail("test@example.com"))
	assert.Equal(t, "", GetDomainFromEmail("invalid-email"))
}

func TestIsZampEmail(t *testing.T) {
	assert.True(t, IsZampEmail("test@zamp.ai"))
	assert.False(t, IsZampEmail("test@example.com"))
}

func TestGetNameFromEmail(t *testing.T) {
	assert.Equal(t, "Test", GetNameFromEmail("test@example.com"))
	assert.Equal(t, "Test", GetNameFromEmail("test.name@example.com"))
}
