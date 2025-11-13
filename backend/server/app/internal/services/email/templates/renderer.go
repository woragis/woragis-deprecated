package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path"
	"strings"
)

//go:embed templates
var templateFS embed.FS

// TemplateID enumerates supported email templates.
type TemplateID string

const (
	TemplateConfirmation  TemplateID = "confirmation"
	TemplatePasswordReset TemplateID = "password_reset"
	TemplateWelcome       TemplateID = "welcome"
	TemplateSessionAlert  TemplateID = "session_alert"
)

// RenderInput holds data required to build an email.
type RenderInput struct {
	Template TemplateID
	Locale   string
	Data     map[string]any
}

// RenderOutput contains the generated email.
type RenderOutput struct {
	Subject string
	HTML    string
	Text    string
}

// Renderer builds localized transactional e-mails.
type Renderer struct {
	defaultLocale string
	funcs         template.FuncMap
}

// NewRenderer constructs a template renderer.
func NewRenderer(defaultLocale string) *Renderer {
	if defaultLocale == "" {
		defaultLocale = "en"
	}

	funcs := template.FuncMap{
		"title": strings.Title,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	return &Renderer{
		defaultLocale: defaultLocale,
		funcs:         funcs,
	}
}

// Render executes the template resolving locale fallback.
func (r *Renderer) Render(input RenderInput) (RenderOutput, error) {
	if input.Template == "" {
		return RenderOutput{}, fmt.Errorf("template id required")
	}

	locale := input.Locale
	if locale == "" {
		locale = r.defaultLocale
	}

	templatePaths := []string{
		path.Join("templates", locale, string(input.Template)+".html"),
		path.Join("templates", r.defaultLocale, string(input.Template)+".html"),
	}

	for _, templatePath := range templatePaths {
		rendered, err := r.renderWithPath(templatePath, input.Data)
		if err == nil {
			return rendered, nil
		}

		// Try next template only when file not found.
		if !isNotFound(err) {
			return RenderOutput{}, err
		}
	}

	return RenderOutput{}, fmt.Errorf("template %s not found for locale %s", input.Template, locale)
}

// Templates returns list of available templates.
func (r *Renderer) Templates() []TemplateID {
	return []TemplateID{
		TemplateConfirmation,
		TemplatePasswordReset,
		TemplateWelcome,
		TemplateSessionAlert,
	}
}

// Preview renders all templates with provided payload for testing.
func (r *Renderer) Preview(locale string, data map[string]any) (map[TemplateID]RenderOutput, error) {
	results := make(map[TemplateID]RenderOutput)
	for _, id := range r.Templates() {
		output, err := r.Render(RenderInput{
			Template: id,
			Locale:   locale,
			Data:     data,
		})
		if err != nil {
			return nil, err
		}
		results[id] = output
	}
	return results, nil
}

func (r *Renderer) renderWithPath(templatePath string, data map[string]any) (RenderOutput, error) {
	files := []string{
		"templates/layout.html",
		templatePath,
	}

	tmpl, err := template.New("email").Funcs(r.funcs).ParseFS(templateFS, files...)
	if err != nil {
		return RenderOutput{}, err
	}

	var htmlBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&htmlBuf, "layout", data); err != nil {
		return RenderOutput{}, err
	}

	var subjectBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&subjectBuf, "subject", data); err != nil {
		return RenderOutput{}, err
	}

	var textBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&textBuf, "text", data); err != nil {
		return RenderOutput{}, err
	}

	return RenderOutput{
		Subject: strings.TrimSpace(subjectBuf.String()),
		HTML:    htmlBuf.String(),
		Text:    strings.TrimSpace(textBuf.String()),
	}, nil
}

func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "file does not exist")
}
