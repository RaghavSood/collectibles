package templates

import (
	"embed"
	"html/template"
	"io"
	"time"

	"github.com/RaghavSood/collectibles/util"
)

//go:embed *
var Templates embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	funcMap := template.FuncMap{
		"Now":               time.Now,
		"NowSecond":         func() time.Time { return time.Now().Truncate(time.Second) },
		"NoEscape":          util.NoEscapeHTML,
		"BTCValueToUSD":     util.BTCValueToUSD,
		"FormatNumber":      util.FormatNumber,
		"MultiParam":        util.MultiParam,
		"ItemPercentString": util.ItemPercentString,
		"ShortUTCTime":      util.ShortUTCTime,
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(Templates, "footer.tmpl", "base.tmpl", "header.tmpl", "series_card.tmpl", "notes.tmpl", "address_list.tmpl", "embed_base.tmpl", "masterlist_table.tmpl"))
	return &Template{
		templates: templates,
	}
}

func RenderSingle(w io.Writer, filename string, data interface{}) error {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{}).ParseFS(Templates, filename))

	return tmpl.ExecuteTemplate(w, filename, data)
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

func (t *Template) RenderNonBase(w io.Writer, newBase, contentTemplate string, data interface{}) error {
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}

	// Parse the specific content template
	_, err = tmpl.ParseFS(Templates, contentTemplate)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, newBase, data)
}
