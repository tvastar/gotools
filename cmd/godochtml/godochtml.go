// Command godochtml converts any packages it is given into go
// documentation.
package main

import (
	"flag"
	"fmt"
	"go/doc"
	"os"

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
	pkgs, err := packages.Load(&packages.Config{Mode: mode}, flag.Args()...)
	if err != nil {
		return err
	}

	if n := packages.PrintErrors(pkgs); n > 0 {
		return fmt.Errorf("%d errors", n)
	}

	for _, pkg := range pkgs {
		p, err := doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath)
		if err != nil {
			return err
		}

		fmt.Println("File: " + p.ImportPath)
		if err := dochtml.Write(os.Stdout, p); err != nil {
			return err
		}
	}

	return nil
}
