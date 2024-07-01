package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/RaghavSood/collectibles/util"
)

//go:embed *
var Templates embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	funcMap := template.FuncMap{
		"NoEscape":      util.NoEscapeHTML,
		"BTCValueToUSD": util.BTCValueToUSD,
		"FormatNumber":  util.FormatNumber,
		"MultiParam":    util.MultiParam,
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(Templates, "footer.tmpl", "base.tmpl", "header.tmpl", "series_card.tmpl", "notes.tmpl", "address_list.tmpl"))
	return &Template{
		templates: templates,
	}
}

func (t *Template) Render(w io.Writer, contentTemplate string, data interface{}) error {
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}

	// Parse the specific content template
	_, err = tmpl.ParseFS(Templates, contentTemplate)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "base.tmpl", data)
}
