// Package dochtml generates html for godochtml.
package dochtml

import (
	"go/doc"
	"html/template"
	"io"

	"github.com/Masterminds/sprig"
)

// Write generataes the html for a specific package.
func Write(w io.Writer, p *doc.Package) error { //nolint: funlen
	fns := map[string]interface{}{
		"sprig": func() interface{} { return sprig.FuncMap() },
	}
	exec := func(t *template.Template, err error) error {
		if err != nil {
			return err
		}

		return t.Funcs(fns).Execute(w, p)
	}

	return exec(template.New("html").Funcs(fns).Parse(`
{{ $desc := index (call sprig.splitList "\n" .Doc) 0 }}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="description" content="{{ .Doc }}">
    <title>{{ .Name }} - {{ $desc }}</title>
  </head>
  <body>
    <h1>Package {{ .Name }}</h1>
    <div id="toc">
      <dl><dd><code>import "{{ .ImportPath }}"</code></dd></dl>
      <dl><dd><a href="#overview" class="overviewLink">Overview</a></dd></dl>
    </div>
    <div id="overview">
      <h2 class="toggle" title="Click to hide Overview section">Overview ▾</h2>
      <p>{{ .Doc }}</p>
    </div>
  </body>
</html>
`))
}
