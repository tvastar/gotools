// Command goinstall installs golang binaries.
//
// Goinstall installs any binaries within a golang module.
//
// The "go get -u" command or "go install" command both have several
// issues when used with modules: they modify the go.mod in the calling
// workspace, they do not apply the go.mod of the golang module being
// installed.  Goinstall works around those limitations by running "go
// get" on an empty module and using that as a the basis.
//
// Usage:
//
//    goinstall [options] mod_path
//
// The mod path can be any path allowed by "go get" though it should
// not contain trailing "/..." -- that is implied.
//
// Options:
//
//   -go string
//     	go binary full path
//   -h	help
//   -o string install directory (default "./bin")
//
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tvastar/gotools/pkg/install"
)

func main() {
	output := flag.String("o", "./bin", "install directory")
	gobin := flag.String("go", "", "go binary full path")
	h := flag.Bool("h", false, "help")

	flag.CommandLine.Usage = usage
	flag.Parse()

	if *h {
		help()
		return
	}

	if importPath := flag.Arg(0); importPath == "" {
		usage()
	} else if err := install.All(*gobin, importPath, *output); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}

func help() {
	fmt.Print(`
Goinstall installs any binaries within a golang module.

The "go get -u" command or "go install" command both have several
issues when used with modules: they modify the go.mod in the calling
workspace, they do not apply the go.mod of the golang module being
installed.  Goinstall works around those limitations by running "go
get" on an empty module and using that as a the basis.

`)
	usage()
}

func usage() {
	fmt.Print(`
Usage:

   goinstall [options] mod_path

The mod path can be any path allowed by "go get" though it should
not contain trailing "/..." -- that is implied.

Options:

`)
	flag.PrintDefaults()
}
