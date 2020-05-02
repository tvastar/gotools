// Command godochtml converts any packages it is given into go
// documentation.
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
	flag.Parse()
	if err := docgen(); err != nil {
		fmt.Fprintf(os.Stderr, "godochtml: %v\n", err)
		os.Exit(1)
	}
}

func docgen() error {
	mode := packages.NeedName | packages.NeedFiles |
		packages.NeedCompiledGoFiles | packages.NeedDeps |
		packages.NeedImports | packages.NeedTypes |
		packages.NeedTypesSizes | packages.NeedSyntax |
		packages.NeedTypesInfo
	cfg := &packages.Config{Mode: mode, Tests: true, Fset: token.NewFileSet()}
	pkgs, err := packages.Load(cfg, flag.Args()...)
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

		if err := dochtml.Write(os.Stdout, p, cfg.Fset); err != nil {
			return err
		}
	}

	return nil
}
