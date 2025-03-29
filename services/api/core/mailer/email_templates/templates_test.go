package emailtemplates

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplate(t *testing.T) {
	// Create temp dir for test templates
	tmpDir, err := os.MkdirTemp("", "email_templates_*")
	assert.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	// Create templater instance
	templater := NewTemplater(tmpDir)

	// Test cases
	tests := []struct {
		name        string
		templateID  string
		template    string
		data        interface{}
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:       "basic template",
			templateID: "test1",
			template:   "Hello {{.Name}}!",
			data:       map[string]string{"Name": "World"},
			want:       "Hello World!",
		},
		{
			name:       "template with multiple variables",
			templateID: "test2",
			template:   "{{.Greeting}} {{.Name}}! How are you {{.Time}}?",
			data: map[string]string{
				"Greeting": "Hi",
				"Name":     "Bob",
				"Time":     "today",
			},
			want: "Hi Bob! How are you today?",
		},
		{
			name:        "non-existent template",
			templateID:  "nonexistent",
			wantErr:     true,
			errContains: "failed to parse template",
		},
		{
			name:        "invalid template syntax",
			templateID:  "invalid",
			template:    "Hello {{.Name!}",
			data:        map[string]string{"Name": "World"},
			wantErr:     true,
			errContains: "failed to parse template",
		},
		{
			name:       "missing data field",
			templateID: "missing",
			template:   "Hello {{.Name}}!",
			data:       map[string]string{"Wrong": "World"},
			want:       "Hello !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create template file if template content provided
			if tt.template != "" {
				templatePath := filepath.Join(tmpDir, tt.templateID+".html")
				err := os.WriteFile(templatePath, []byte(tt.template), 0644)
				assert.NoError(t, err, "Failed to write template file")
			}

			// Call GetTemplate through the templater instance
			got, err := templater.GetTemplate(tt.templateID, tmpDir, tt.data)
			fmt.Println("got", got)
			fmt.Println("err", err)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
