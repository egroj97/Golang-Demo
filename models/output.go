package models

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// DataEntryJSON will be used to return properly formatted information on rest
// endpoint
type DataEntryJSON struct {
	*DataEntry
	Keywords     []string `json:"keywords"`
	BureauCodes  []string `json:"bureauCodes"`
	ProgramCodes []string `json:"programCodes"`
}

// PageData contains all the data sent to the template to be rendered.
type PageData struct {
	Payloads []Payload
}

// Template implemented the Renderer interface from echo.
type Template struct {
	Templates *template.Template
}

// GetCreatedAtHuman is an aux. method to be used on the template.
func (p *Payload) GetCreatedAtHuman() string {
	return p.CreatedAt.Format("January 2, 2006	")
}

// Render is the method needed to implement the Renderer interface.
func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	// In this case the html/template template engine is used.
	return t.Templates.ExecuteTemplate(w, name, data)
}
