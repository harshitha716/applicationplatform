package templates

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTemplateLoader(t *testing.T) {
	loader := InitTemplateLoader("server/routes/admin/templates")
	assert.NotNil(t, loader)
	assert.IsType(t, &templateLoader{}, loader)
}

func TestLoadTemplate(t *testing.T) {
	tests := []struct {
		name         string
		templateName string
		setup        func() (string, func())
		data         TemplateData

		expectedError bool
	}{
		{
			name:         "Load existing template",
			templateName: "admin-home.html",
			setup: func() (string, func()) {
				// create temp html file and return path
				tempFile, err := os.CreateTemp("", "admin-home.html")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				defer tempFile.Close()
				return tempFile.Name(), func() {
					os.Remove(tempFile.Name())
				}
			},
			data: TemplateData{
				Title:       "Test Title",
				Environment: "test",
				Route:       "/test",
				Data:        struct{}{},
			},
			expectedError: false,
		},
		{
			name:         "Load non-existent template",
			templateName: "non-existent.html",
			setup: func() (string, func()) {
				return "", func() {}
			},
			data: TemplateData{
				Title:       "Test Title",
				Environment: "test",
				Route:       "/test",
				Data:        struct{}{},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			filename, cleanup := tt.setup()

			testMap := map[string]string{
				"admin-home.html": filename,
			}
			loader := initTemplateLoader(testMap)

			tmpl, err := loader.LoadTemplate(tt.templateName, tt.data)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, tmpl)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tmpl)
			}
			cleanup()
		})
	}
}

func TestExecuteTemplate(t *testing.T) {
	tests := []struct {
		name          string
		templateName  string
		setup         func() (string, func())
		data          TemplateData
		expectedError bool
	}{
		{
			name:         "Execute existing template",
			templateName: "admin-home.html",
			setup: func() (string, func()) {
				return "admin-home.html", func() {

				}
			},
			data: TemplateData{
				Title:       "Test Title",
				Environment: "test",
				Route:       "/test",
				Data:        struct{}{},
			},
			expectedError: false,
		},
		{
			name:         "Execute non-existent template",
			templateName: "non-existent.html",
			setup: func() (string, func()) {
				return "", func() {}
			},
			data: TemplateData{
				Title:       "Test Title",
				Environment: "test",
				Route:       "/test",
				Data:        struct{}{},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, cleanup := tt.setup()
			testMap := map[string]string{
				"admin-home.html": filename,
			}
			loader := initTemplateLoader(testMap)
			buf := new(bytes.Buffer)
			err := loader.ExecuteTemplate(buf, tt.templateName, tt.data)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, buf.String())
			}
			cleanup()
		})
	}
}
