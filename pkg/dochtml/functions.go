package dochtml

import (
	"bytes"
	"go/doc"
	"go/format"
	"go/token"
	"html/template"

	"github.com/Masterminds/sprig"
)

// Functions builds out the standard functions needed by the html template.
type Functions struct {
	*doc.Package
	*token.FileSet
}

// Map returns the default function map
func (f *Functions) Map() map[string]interface{} {
	sprigf := sprig.FuncMap()
	astf := f.astFunctions()
	docf := f.docFunctions()
	unsafef := f.unsafeFunctions()
	return map[string]interface{}{
		"ast":         func() interface{} { return astf },
		"doc":         func() interface{} { return docf },
		"sprig":       func() interface{} { return sprigf },
		"unsafe":      func() interface{} { return unsafef },
		"highlighter": func() interface{} { return &Highlighter{} },
	}
}

func (f *Functions) astFunctions() map[string]interface{} {
	highlight := &Highlighter{}
	format := func(node interface{}) (string, error) {
		var buf bytes.Buffer
		if err := format.Node(&buf, f.FileSet, node); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	return map[string]interface{}{
		"text": func(node interface{}) (string, error) {
			return format(node)
		},
		"html": func(node interface{}) (template.HTML, error) {
			s, err := format(node)
			return template.HTML(highlight.String(s)), err //nolint: gosec
		},
	}
}

func (f *Functions) docFunctions() map[string]interface{} {
	highlight := &Highlighter{}
	return map[string]interface{}{
		"synopsis": doc.Synopsis,
		"toHTML": func(s string) string {
			var buf bytes.Buffer
			doc.ToHTML(&buf, s, nil)
			return highlight.HTML(buf.String())
		},
	}
}

func (f *Functions) unsafeFunctions() map[string]interface{} {
	return map[string]interface{}{
		"html": func(s string) interface{} { return template.HTML(s) }, //nolint: gosec
		"js":   func(s string) interface{} { return template.JS(s) },   //nolint: gosec
		"css":  func(s string) interface{} { return template.CSS(s) },
		"url":  func(s string) interface{} { return template.URL(s) },      //nolint: gosec
		"attr": func(s string) interface{} { return template.HTMLAttr(s) }, //nolint: gosec
	}
}
