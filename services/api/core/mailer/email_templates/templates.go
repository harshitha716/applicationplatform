package emailtemplates

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

type Templater interface {
	GetTemplate(templateID string, templatesPath string, data interface{}) (string, error)
}

type templater struct {
	emailTemplatesPath string
}

func NewTemplater(emailTemplatesPath string) Templater {
	return templater{emailTemplatesPath: emailTemplatesPath}
}

// GetTemplate loads and executes an email template with the given data
func (t templater) GetTemplate(templateID string, templatesPath string, data interface{}) (string, error) {
	// Build template path
	templatePath := filepath.Join(templatesPath, fmt.Sprintf("%s.html", templateID))

	// Parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateID, err)
	}

	// Execute template into string builder
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateID, err)
	}

	return result.String(), nil
}
