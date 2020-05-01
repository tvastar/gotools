// Package install implements installing golang tools predictably.
package install

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/semver"
)

// All installs all the binaries for a versioned import path.
//
// gobin is the location of the go binary.
//
// importPath is the import path to install.
//
// dir is the directory where the binaries should be installed. It is
// created if needed.
func All(gobin, importPath, dir string) error {
	var err error

	// Convert dir and gobin to absolute paths.
	if dir, err = resolveOutputDir(dir); err != nil {
		return err
	}
	if gobin, err = resolveGobinPath(gobin); err != nil {
		return err
	}

	// Create temporary directory and ensure it is cleaned up.
	tmp, err := ioutil.TempDir("", "goinstall")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	return run(tmp,
		exec.Command(gobin, "mod", "init", "x"),
		exec.Command(gobin, "get", importPath),
		exec.Command(gobin, "build", "-o", dir, allbins(importPath)),
	)
}

func run(dir string, commands ...*exec.Cmd) error {
	for _, cmd := range commands {
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Println(cmd.Args, err, "\n"+string(out))
			return err
		}
	}
	return nil
}

func allbins(importPath string) string {
	if idx := strings.LastIndex(importPath, "@"); idx != -1 && semver.IsValid(importPath[idx+1:]) {
		return importPath[:idx] + "/..."
	}
	return importPath + "/..."
}

func resolveGobinPath(gobin string) (string, error) {
	if gobin == "" {
		exe := "go"
		if runtime.GOOS == "windows" {
			exe = "go.exe"
		}

		return exec.LookPath(exe)
	}
	return filepath.Abs(gobin)
}

func resolveOutputDir(dir string) (string, error) {
	var err error

	if dir, err = filepath.Abs(dir); err != nil {
		return dir, err
	}

	return dir, os.MkdirAll(dir, 0766)
}
