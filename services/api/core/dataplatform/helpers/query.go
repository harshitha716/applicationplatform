package helpers

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"go.uber.org/zap"
)

func FillQueryTemplate(ctx context.Context, query string, params map[string]string) (string, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	sanitizedKeys := make(map[string]string) // Map to track original and sanitized keys
	query = replaceTemplateKeys(query, sanitizedKeys)

	// Step 2: Update params only for the sanitized keys
	sanitizedParams := make(map[string]string)
	for originalKey, sanitizedKey := range sanitizedKeys {
		if value, exists := params[originalKey]; exists {
			sanitizedParams[sanitizedKey] = value
		}
	}

	template, err := template.New("query").Parse(query)
	if err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return "", errors.ErrTemplateParsingFailed
	}

	var filledQuery bytes.Buffer
	if err := template.Execute(&filledQuery, sanitizedParams); err != nil {
		logger.Error(errors.TemplateParsingFailedErrMessage, zap.Error(err))
		return "", errors.ErrTemplateParsingFailed
	}

	filledQueryString := filledQuery.String()
	return filledQueryString, nil
}

func replaceTemplateKeys(query string, sanitizedKeys map[string]string) string {
	re := regexp.MustCompile(`{{\.\s*([\w-]+)\s*}}`)
	return re.ReplaceAllStringFunc(query, func(match string) string {
		// Extract the key inside the {{.}} and sanitize it
		key := strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{.")
		sanitizedKey := strings.ReplaceAll(key, "-", "_")
		sanitizedKeys[key] = sanitizedKey // Track the original and sanitized keys
		return "{{." + sanitizedKey + "}}"
	})
}

func AddCommentsToQuery(query string, metadata map[string]string) string {
	comment := ""
	for key, value := range metadata {
		comment += fmt.Sprintf("\n-- %s='%s'", key, value)
	}
	return query + comment
}

func BuildDatabricksTableName(catalog string, schema string, table string) string {
	return fmt.Sprintf("`%s`.`%s`.`%s`", catalog, schema, table)
}
