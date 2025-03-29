package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidShortInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "string length 1",
			input: "a",
			want:  true,
		},
		{
			name:  "string length 24",
			input: "123456789012345678901234",
			want:  true,
		},
		{
			name:  "string length 25",
			input: "1234567890123456789012345",
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsValidShortInput(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidMediumInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "string length 1",
			input: "a",
			want:  true,
		},
		{
			name:  "string length 64",
			input: "1234567890123456789012345678901234567890123456789012345678901234",
			want:  true,
		},
		{
			name:  "string length 65",
			input: "12345678901234567890123456789012345678901234567890123456789012345",
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsValidMediumInput(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidLongInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "string length 1",
			input: "a",
			want:  true,
		},
		{
			name:  "string length 255",
			input: "234567890123456789456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			want:  true,
		},
		{
			name:  "string length 256",
			input: "1234567895678901234a67890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsValidLongInput(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidHexCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid hex code lowercase",
			input: "#abc123",
			want:  true,
		},
		{
			name:  "valid hex code uppercase",
			input: "#ABC123",
			want:  true,
		},
		{
			name:  "valid hex code mixed case",
			input: "#aBC123",
			want:  true,
		},
		{
			name:  "missing hash prefix",
			input: "abc123",
			want:  false,
		},
		{
			name:  "too short",
			input: "#abc12",
			want:  false,
		},
		{
			name:  "too long",
			input: "#abc1234",
			want:  false,
		},
		{
			name:  "invalid characters",
			input: "#abc12g",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "only hash",
			input: "#",
			want:  false,
		},
		{
			name:  "valid hex code with numbers",
			input: "#000000",
			want:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsValidHexCode(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
