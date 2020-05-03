//nolint: lll
// Command godochtml generates package documentation as HTML.
//
// The output is very close to (but not exactly same as) godoc. The
// generated HTML is simple static HTML but includes syntax
// highlighting.
//
// Usage: godochtml [options] "package or package directory"
//
//    -src fileregexp=url -- use this to make types and functions
//  links.
//
// Example:
//
//     godochtml -src 'github.com/([^/]*)/([^/]*)(.*)=https://github.com/$1/$2/blob/master$3' github.com/tvastar/gotools/cmd/godochtml
//
// The generated HTML is written to console.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/tvastar/gotools/pkg/dochtml"
)

func main() {
	l := &dochtml.FileLinker{}

	flag.CommandLine.Usage = help
	h := flag.Bool("h", false, "help")
	flag.Var(l, "src", `source url map

Source urls can be mapped using find_regexp=replacement_pattern syntax.
For example, the following argument maps github import paths to github source urls:
   -src 'github.com/([^/]*)/([^/]*)(.*)=https://github.com/$1/$2/blob/master$3'
`)
	flag.Parse()

	if *h || flag.Arg(0) == "" {
		help()
		return
	}

	if err := docgen(flag.Arg(0), l); err != nil {
		fmt.Fprintf(os.Stderr, "godochtml: %v\n", err)
		os.Exit(1)
	}
}

func docgen(pattern string, linker *dochtml.FileLinker) error {
	mode := packages.NeedName | packages.NeedFiles |
		packages.NeedCompiledGoFiles | packages.NeedDeps |
		packages.NeedImports | packages.NeedTypes |
		packages.NeedTypesSizes | packages.NeedSyntax |
		packages.NeedTypesInfo
	cfg := &packages.Config{Mode: mode, Tests: true, Fset: token.NewFileSet()}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return err
	}

	if n := packages.PrintErrors(pkgs); n > 0 {
		return fmt.Errorf("%d errors", n)
	}

	// filesets tracks all files for a given package path name.
	filesets := map[string][]*ast.File{}
	for _, pkg := range pkgs {
		path := pkg.PkgPath
		if strings.HasSuffix(path, "_test") || strings.HasSuffix(path, ".test") {
			path = path[:len(path)-5]
		}
		for _, file := range pkg.Syntax {
			if f := cfg.Fset.File(file.Pos()); f != nil && strings.HasSuffix(f.Name(), ".go") {
				filesets[path] = append(filesets[path], file)
			}
		}
	}

	for _, files := range filesets {
		p, err := doc.NewFromFiles(cfg.Fset, files, pattern)
		if err != nil {
			return err
		}

		p.Examples = doc.Examples(files...)

		if err := dochtml.Write(os.Stdout, p, cfg.Fset, "", linker); err != nil {
			return err
		}
	}

	return nil
}

func help() {
	fmt.Print(`
Usage: godochtml <pkg_name>

pkg_name can be a local directory or an import path.

The generated HTML is written to console.
`)
	flag.PrintDefaults()
}
