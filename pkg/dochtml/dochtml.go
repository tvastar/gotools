// Package dochtml generates html for godochtml.
package dochtml

import (
	"bytes"
	"go/doc"
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
		"html": func(s string) interface{} { return template.HTML(s) },
		"js":   func(s string) interface{} { return template.JS(s) },
		"css":  func(s string) interface{} { return template.CSS(s) },
		"url":  func(s string) interface{} { return template.URL(s) },
		"attr": func(s string) interface{} { return template.HTMLAttr(s) },
	}
}

// Write generataes the html for a specific package.
func Write(w io.Writer, p *doc.Package) error { //nolint: funlen
	fns := map[string]interface{}{
		"sprig":  func() interface{} { return sprig.FuncMap() },
		"doc":    docFunctions,
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
    <div id="toc">
      <dl><dd><code>import "{{ .ImportPath }}"</code></dd></dl>
      <dl><dd><a href="#overview" class="overviewLink">Overview</a></dd></dl>
    </div>
    <div id="overview">
      <h2 class="toggle" title="Click to hide Overview section">Overview ▾</h2>
      <p>{{ call unsafe.html $overview }}</p>
    </div>
  </body>
</html>
`))
}
