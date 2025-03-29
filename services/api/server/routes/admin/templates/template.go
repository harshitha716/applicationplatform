package templates

import (
	"fmt"
	"html/template"
	"io"

	goform "github.com/joncalhoun/form"
)

type TemplateData struct {
	Title       string
	Environment string
	Route       string
	Data        any
}

type TemplateLoader interface {
	LoadTemplate(templateName string, data TemplateData) (*template.Template, error)
	ExecuteTemplate(w io.Writer, name string, data TemplateData) error
}

type templateLoader struct {
	namePathMap map[string]string
}

func InitTemplateLoader(templatesDir string) TemplateLoader {

	templateNamePathMap := map[string]string{
		"admin-home.html": fmt.Sprintf("%s/admin-home.html", templatesDir),
		"base-form.html":  fmt.Sprintf("%s/base-form.html", templatesDir),
	}

	return initTemplateLoader(templateNamePathMap)
}

func initTemplateLoader(tnMap map[string]string) TemplateLoader {
	return &templateLoader{
		namePathMap: tnMap,
	}
}

func (t *templateLoader) LoadTemplate(templateName string, data TemplateData) (*template.Template, error) {

	formBuilderTemplate := template.Must(template.New(templateName).Parse(`
<div class="mb-4">
	<label class="block text-grey-darker text-sm font-bold mb-2" {{with .ID}}for="{{.}}"{{end}}>
		{{.Label}}
	</label>
	<input class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight" {{with .ID}}id="{{.}}"{{end}} type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
	{{with .Footer}}
		<p class="text-grey pt-2 text-xs italic">{{.}}</p>
	{{end}}
</div>
	`))

	formBuilder := goform.Builder{
		InputTemplate: formBuilderTemplate,
	}

	templatePath, ok := t.namePathMap[templateName]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	newTemplate := template.New(templateName)
	newTemplate.Funcs(formBuilder.FuncMap())
	newTemplate.Funcs(template.FuncMap{
		"pre_process_route": func(_ any) string {
			return data.Route
		},
		"pre_process_environment": func(_ any) string {
			return data.Environment
		},
		"pre_process_title": func(_ any) string {
			return data.Title
		},
	})
	htmlTemplate, err := template.Must(newTemplate, nil).ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %s", err)
	}

	return htmlTemplate, nil
}

func (t *templateLoader) ExecuteTemplate(w io.Writer, name string, data TemplateData) error {

	htmlTemplate, err := t.LoadTemplate(name, data)

	if err != nil {
		return err
	}

	return htmlTemplate.ExecuteTemplate(w, name, data.Data)
}
