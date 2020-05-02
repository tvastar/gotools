// Package dochtml generates html for godochtml.
package dochtml

import (
	"go/doc"
	"go/token"
	"html/template"
	"io"
)

// Write generates the html for a specific package using the provided
// template.  If no template is provided, the default template is
// chosen..
func Write(w io.Writer, p *doc.Package, fset *token.FileSet, tpl string) error {
	if tpl == "" {
		tpl = DefaultTemplate()
	}

	fns := &Functions{p, fset}
	t, err := template.New("html").Funcs(fns.Map()).Parse(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, p)
}

// DefaultTemplate returns the default HTML template used by godochtml.
//nolint: lll, funlen
func DefaultTemplate() string {
	return `
{{- $synopsis := call doc.synopsis .Doc -}}
{{- $overview := call doc.toHTML .Doc -}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="description" content="{{ .Doc }}">
    <title>{{ .Name }} - {{ $synopsis }}</title>
    <style type="text/css">
      html, body { padding: 0; border: 0; margin: 0 }
      body {
        max-width: 728px;
        padding: 0 15px;
        margin: 0 auto;
        font-family: "Helvetica Neue",Helvetica,Arial,sans-serif;
        font-size: 14.5px;
        line-height: 1.42;
        color: #333;
      }
      h2, .h2 { font-size: 30px; }
      h3, .h3 { font-size: 24px; }
      h4, .h4 { font-size: 20px; }
      h1, .h1, h2, .h2, h3, .h3, h4, .h4, h5, .h5 {
        margin: 20px 0 10px 0;
        font-weight: 500;
        line-height: 1.1;
      }
      p { margin: 0 0 10px; }
      pre {
        padding: 9.5px;
        margin: 0 0 10px;
        font-size: 13px;
        line-height: 1.4;
        word-break: break-all;
        word-wrap: break-word;
        background-color: #f5f5f5;
        border: 1px solid #ccc;
        border-radius: 4px;
        box-sizing: border-box;
      }
      code, kbd, pre, samp {
        font-family: Menlo,Monaco,Consolas,"Courier New",monospace;
      }
      #pkg-index > ul, #pkg-examples > ul { padding-left: 0; }
      #pkg-index > ul > li, #pkg-examples > ul > li { list-style: none; }
      #pkg-index > ul > ul > li { list-style-type: circle; }
      a { text-decoration: none; color: #375eab; }
      details > summary { color: #375eab; margin: 10px 0; }

      {{ highlighter.Styles }}
    </style>
  </head>
  <body>
    <h1>Package {{ .Name }}</h1>
    <div id="pkg-overview">
      <p>{{ call unsafe.html $overview }}</p>
    </div>
    <div id="pkg-index">
      <h3>Index</h2>
      <ul>
      {{ if .Consts }}<li><a href="#pkg-consts">Constants</a></li>{{ end }}
      {{ if .Vars }}<li><a href="#pkg-vars">Variables</a></li>{{ end }}
      {{ range .Funcs }}<li><a href="#{{.Name}}">{{ call ast.text .Decl }}</a></li>
      {{ end }}{{ range .Types }}<li><a href="#{{.Name}}">type {{ .Name }}</a></li>
        {{ if .Funcs }}<ul>{{ range .Funcs }}<li><a href="#{{.Name}}">{{ call ast.text .Decl }}</a></li>{{ end }}</ul>{{ end }}
        {{ if .Methods }}<ul>{{ range .Methods }}<li><a href="#{{.Name}}">{{ call ast.text .Decl }}</a></li>{{ end }}</ul>{{ end }}
      {{ end }}
      </ul>
    </div>
    <div id="pkg-examples">
    {{ if .Examples }}<h4>Examples</h4>
      <ul>
        {{ range .Examples }}<li><a href="#{{.Name}}">{{ .Name }}</a></li>{{ end }}
      </ul>
    {{ end }}
    </div>
    <div id="pkg-consts">
    {{ if .Consts}}<h4>Constants</h4>
       {{ range .Consts }}
         {{ call ast.html .Decl }}
         <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
       {{ end }}
    {{ end }}
    </div>
    <div id="pkg-vars">
    {{ if .Vars}}<h4>Variables</h4>
       {{ range .Vars }}
         {{ call ast.html .Decl }}
         <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
       {{ end }}
    {{ end }}
    </div>
    {{ range .Funcs }}<h3 id="{{.Name}}">func {{.Name}}</h3>
       {{ call ast.html .Decl }}
       <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
       {{ range .Examples }}
          <details>
            {{ if .Name }}<summary>Example ({{.Name}})</summary>
            {{ else }}<summary>Example</summary>
            {{ end }}
            <p>Code:</p>
            {{ call ast.html .Code }}
            {{ if .Output }}<p>Output:</p><pre>{{ .Output }}</pre>{{ end }}
            <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
          </details>
       {{ end }}
    {{ end }}
    {{ range .Types }}<h3 id="{{.Name}}">type {{.Name}}</h3>
       {{ call ast.html .Decl }}
       <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
       {{ range .Examples }}
          <details>
            {{ if .Name }}<summary>Example ({{.Name}})</summary>
            {{ else }}<summary>Example</summary>
            {{ end }}
            <p>Code:</p>
            {{ call ast.html .Code }}
            {{ if .Output }}<p>Output:</p><pre>{{ .Output }}</pre>{{ end }}
            <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
          </details>
       {{ end }}
       {{ range .Consts }}
         {{ call ast.html .Decl }}
         <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
       {{ end }}
       {{ range .Vars }}
         {{ call ast.html .Decl }}
         <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
       {{ end }}
       {{ range .Funcs }}<h3 id="{{.Name}}">func {{if .Recv}}({{.Recv}}){{end}} {{.Name}}</h3>
         {{ call ast.html .Decl }}
         <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
         {{ range .Examples }}
            <details>
              {{ if .Name }}<summary>Example ({{.Name}})</summary>
              {{ else }}<summary>Example</summary>
              {{ end }}
              <p>Code:</p>
              {{ call ast.html .Code }}
              {{ if .Output }}<p>Output:</p><pre>{{ .Output }}</pre>{{ end }}
              <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
            </details>
         {{ end }}
       {{ end }}
       {{ range .Methods }}
         <h3 id="{{.Name}}">func {{if .Recv}}({{.Recv}}){{end}} {{.Name}}</h3>
         {{ call ast.html .Decl }}
         <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
         {{ range .Examples }}
            <details>
              {{ if .Name }}<summary>Example ({{.Name}})</summary>
              {{ else }}<summary>Example</summary>
              {{ end }}
              <p>Code:</p>
              {{ call ast.html .Code }}
              {{ if .Output }}<p>Output:</p><pre>{{ .Output }}</pre>{{ end }}
              <p>{{ call unsafe.html (call doc.toHTML .Doc)}}</p>
            </details>
         {{ end }}
       {{ end }}
    {{ end }}
  </body>
</html>
`
}
