package tests

import (
	"testing"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	"github.com/stretchr/testify/assert"
)

func TestConvertToJSONString(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    string
		expectError bool
	}{
		{
			name:        "Valid input - simple struct",
			input:       struct{ Name string }{Name: "John"},
			expected:    `{"Name":"John"}`,
			expectError: false,
		},
		{
			name:        "Valid input - map",
			input:       map[string]interface{}{"age": 30, "city": "New York"},
			expected:    `{"age":30,"city":"New York"}`,
			expectError: false,
		},
		{
			name:        "Invalid input - channel",
			input:       make(chan int),
			expected:    "",
			expectError: true,
		},
		{
			name:        "Valid input - slice",
			input:       []string{"apple", "banana", "cherry"},
			expected:    `["apple","banana","cherry"]`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helpers.ConvertToJSONString(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expected, result)
			}
		})
	}
}

func TestConvertToJSONStringWithReplacements(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		expected     string
		expectError  bool
		replacements map[string]string
	}{
		{
			name:        "Valid input - simple struct",
			input:       struct{ Name string }{Name: "John"},
			expected:    `{"Name":"Jane"}`,
			expectError: false,
			replacements: map[string]string{
				"John": "Jane",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helpers.ConvertToJSONStringWithReplacements(tt.input, tt.replacements)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expected, result)
			}
		})
	}

}
