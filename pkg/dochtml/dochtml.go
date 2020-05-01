// Package dochtml generates html for godochtml.
package dochtml

import (
	"bytes"
	"go/doc"
	"go/format"
	"go/token"
	"html/template"
	"io"

	"github.com/Masterminds/sprig"
)

func docFunctions() map[string]interface{} {
	return map[string]interface{}{
		"synopsis": doc.Synopsis,
		"toHTML": func(s string) string {
			var buf bytes.Buffer
			doc.ToHTML(&buf, s, nil)
			return buf.String()
		},
	}
}

func unsafeFunctions() map[string]interface{} {
	return map[string]interface{}{
		"html": func(s string) interface{} { return template.HTML(s) }, //nolint: gosec
		"js":   func(s string) interface{} { return template.JS(s) },   //nolint: gosec
		"css":  func(s string) interface{} { return template.CSS(s) },
		"url":  func(s string) interface{} { return template.URL(s) },      //nolint: gosec
		"attr": func(s string) interface{} { return template.HTMLAttr(s) }, //nolint: gosec
	}
}

func astFunctions() map[string]interface{} {
	return map[string]interface{}{
		"text": func(node interface{}) string {
			var buf bytes.Buffer
			if err := format.Node(&buf, token.NewFileSet(), node); err != nil {
				return err.Error()
			}
			return buf.String()
		},
	}
}

// Write generataes the html for a specific package.
func Write(w io.Writer, p *doc.Package) error { //nolint: funlen
	fns := map[string]interface{}{
		"ast":    astFunctions,
		"doc":    docFunctions,
		"sprig":  func() interface{} { return sprig.FuncMap() },
		"unsafe": unsafeFunctions,
	}
	exec := func(t *template.Template, err error) error {
		if err != nil {
			return err
		}

		return t.Funcs(fns).Execute(w, p)
	}

	return exec(template.New("html").Funcs(fns).Parse(`
{{- $synopsis := call doc.synopsis .Doc -}}
{{- $overview := call doc.toHTML .Doc -}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="description" content="{{ .Doc }}">
    <title>{{ .Name }} - {{ $synopsis }}</title>
  </head>
  <body>
    <h1>Package {{ .Name }}</h1>
    <div id="pkg-overview">
      <p>{{ call unsafe.html $overview }}</p>
    </div>
    <div id="pkg-index">
      <h2>Index</h2>
      <ul>
      {{ if .Consts }}<li><a href="#pkg-consts">Constants</a></li>{{ end }}
      {{ if .Vars }}<li><a href="#pkg-vars">Variables</a></li>{{ end }}
      {{ range .Funcs }}<li><a href="#{{.Name}}">{{ call ast.text .Decl }}</a></li>
      {{ end }}{{ range .Types }}<li><a href="#{{.Name}}">type {{ .Name }}</a></li>
      {{ end }}
      </ul>
    </div>
  </body>
</html>
`))
}
