// Package dochtml generates html for godochtml.
package dochtml

import (
	"bytes"
	"go/doc"
	"go/format"
	"go/token"
	"html/template"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/Masterminds/sprig"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func highlightCode(s string) string {
	lexer := lexers.Get("go")
	style, err := styles.Get("arduino").Builder().
		Add(chroma.Background, "#333 bg:#ffffff00").Build()
	if err != nil {
		return s
	}
	formatter := formatters.Get("html")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	iter, err := lexer.Tokenise(nil, s)
	if err != nil {
		return s
	}
	var buf bytes.Buffer
	if err = formatter.Format(&buf, style, iter); err != nil {
		return s
	}
	return buf.String()
}

func highlightPre(s string) string {
	root := &html.Node{Type: html.ElementNode, DataAtom: atom.Div, Data: "div"}
	nn, err := html.ParseFragment(strings.NewReader(s), root)
	if err != nil {
		return s
	}
	for _, n := range nn {
		root.AppendChild(n)
	}
	var visit func(*html.Node)
	visit = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
		if n.DataAtom != atom.Pre {
			return
		}
		data := n.FirstChild.Data
		nn, err = html.ParseFragment(strings.NewReader(highlightCode(data)), nil)
		if err != nil {
			return
		}
		n.Data = "div"
		n.DataAtom = atom.Div
		n.RemoveChild(n.FirstChild)
		for _, nx := range nn {
			n.AppendChild(nx)
		}
	}
	visit(root)
	var result bytes.Buffer
	if err := html.Render(&result, root); err != nil {
		return s
	}
	return result.String()
}

func docFunctions() map[string]interface{} {
	return map[string]interface{}{
		"synopsis": doc.Synopsis,
		"toHTML": func(s string) string {
			var buf bytes.Buffer
			doc.ToHTML(&buf, s, nil)
			return highlightPre(buf.String())
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

func astFunctions(fset *token.FileSet) func() map[string]interface{} {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"text": func(node interface{}) string {
				var buf bytes.Buffer
				if err := format.Node(&buf, fset, node); err != nil {
					return err.Error()
				}
				return buf.String()
			},
			"html": func(node interface{}) interface{} {
				var buf bytes.Buffer
				if err := format.Node(&buf, fset, node); err != nil {
					return err
				}
				return template.HTML(highlightCode(buf.String())) //nolint: gosec
			},
		}
	}
}

// Write generataes the html for a specific package.
//nolint: lll, funlen
func Write(w io.Writer, p *doc.Package, fset *token.FileSet) error {
	fns := map[string]interface{}{
		"ast":    astFunctions(fset),
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
  </body>
</html>
`))
}
