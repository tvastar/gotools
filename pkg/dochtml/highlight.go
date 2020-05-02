package dochtml

import (
	"bytes"
	"strings"

	"html/template"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// Highlighter implements syntax highlighting.
type Highlighter struct {
}

// Styles returns a <style> node with appropriate classes
func (h *Highlighter) Styles() (template.CSS, error) {
	style, err := styles.Get("arduino").Builder().
		Add(chroma.Background, "#333 bg:#ffffff00").
		Build()
	if err != nil {
		return template.CSS(""), nil //nolint: gosec
	}
	formatter := chromahtml.New(chromahtml.WithClasses(true))

	var buf bytes.Buffer
	err = formatter.WriteCSS(&buf, style)
	return template.CSS(buf.String()), err //nolint: gosec
}

// Code converts a static into highlighted HTML.
func (h *Highlighter) String(s string) string {
	lexer := lexers.Get("go")
	style, err := styles.Get("arduino").Builder().
		Add(chroma.Background, "#333 bg:#ffffff00").
		Build()
	if err != nil {
		return s
	}
	formatter := chromahtml.New(chromahtml.WithClasses(true))
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

// HTML scans the HTML for `pre` tags, converting these to highlighted
// HTML.
func (h *Highlighter) HTML(s string) string {
	root := &html.Node{Type: html.ElementNode}
	if err := h.parseFragment(s, root); err != nil {
		return s
	}

	var visit func(*html.Node)
	visit = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
		if n.DataAtom != atom.Pre || !h.isLastTextNode(n.FirstChild) {
			return
		}
		if err := h.parseFragment(h.String(n.FirstChild.Data), n); err != nil {
			return
		}

		n.Data = "div"
		n.DataAtom = atom.Div
		n.RemoveChild(n.FirstChild)
	}
	visit(root)

	if contents, err := h.innerHTML(root); err == nil {
		return contents
	}
	return s
}

func (h *Highlighter) isLastTextNode(n *html.Node) bool {
	return n != nil && n.NextSibling == nil && n.Type == html.TextNode
}

func (h *Highlighter) parseFragment(s string, parent *html.Node) error {
	nn, err := html.ParseFragment(strings.NewReader(s), parent)
	if err != nil {
		return err
	}
	for _, n := range nn {
		parent.AppendChild(n)
	}
	return nil
}

func (h *Highlighter) innerHTML(n *html.Node) (string, error) {
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&buf, c); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
