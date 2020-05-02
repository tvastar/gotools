// Command godochtml generates package documentation as HTML.
//
// The output is very close to (but not exactly same as) godoc. The
// generated HTML is simple static HTML but includes syntax
// highlighting.
//
// Usage: godochtml "package or package directory"
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
	flag.CommandLine.Usage = help
	h := flag.Bool("h", false, "help")
	flag.Parse()

	if *h || flag.Arg(0) == "" {
		help()
		return
	}

	if err := docgen(flag.Arg(0)); err != nil {
		fmt.Fprintf(os.Stderr, "godochtml: %v\n", err)
		os.Exit(1)
	}
}

func docgen(pattern string) error {
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

	for pkgPath, files := range filesets {
		p, err := doc.NewFromFiles(cfg.Fset, files, pkgPath)
		if err != nil {
			return err
		}

		p.Examples = doc.Examples(files...)

		if err := dochtml.Write(os.Stdout, p, cfg.Fset, ""); err != nil {
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
}
