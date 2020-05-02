// Command tag creates a semver tag in the current git repo and pushes
// it.
//
// The semver versioning scheme used is something like: v5.3.2
// (ie vMajor.Minor.Patch).
//
// Tag reads the current git repository for the highest semver and
// based on the command line argument passed, it increments the patch,
// minor or major version.
//
// An exact tag can also be provided for use with the very first
// semver version.
//
// If no command is provided, the current highest semver tag is
// printed.
//
// Usage:
//
// tag [options] (patch|minor|major|"exact-tag")
// commands:
//   patch -- increments the patch version.
//   major -- increments the major version.
//   minor -- increments the minor version.
//   tag   -- use this exact value.
// options:
//   -q       -- do not prompt for confirmation.
//   -m <msg> -- the message to use with the tag.
//   -h       -- help
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"golang.org/x/mod/semver"
)

func main() {
	must := func(err error) {
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}

	flag.CommandLine.Usage = usage
	q := flag.Bool("q", false, "do not prompt for confirmation.")
	m := flag.String("m", "", "the message to use with the tag")
	h := flag.Bool("h", false, "help")
	flag.Parse()

	if *h {
		help()
		return
	}

	r, err := git.PlainOpen(gitDir())
	must(err)

	tag, err := latest(r)
	must(err)

	switch flag.Arg(0) {
	case "":
		fmt.Println("Latest semver tag", tag)
	case "patch":
		must(update(r, *q, *m, next(tag, 0, 0, 1)))
	case "minor":
		must(update(r, *q, *m, next(tag, 0, 1, 0)))
	case "major":
		must(update(r, *q, *m, next(tag, 1, 0, 0)))
	default:
		must(update(r, *q, *m, flag.Arg(0)))
	}
}

func latest(r *git.Repository) (string, error) {
	iter, err := r.Tags()
	if err != nil {
		return "", err
	}

	latest := ""
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		obj, err := r.TagObject(ref.Hash())
		switch err {
		case nil:
			if semver.Compare(obj.Name, latest) > 0 {
				latest = obj.Name
			}
		case plumbing.ErrObjectNotFound:
			err = nil
		}
		return err
	})

	return latest, err
}

func next(current string, major, minor, patch int) string {
	var maj, min, pat int
	if n, err := fmt.Sscanf(current, "v%d.%d.%d", &maj, &min, &pat); n != 3 || err != nil {
		panic("internal error")
	}
	maj += major
	min += minor
	pat += patch
	return fmt.Sprintf("v%d.%d.%d", maj, min, pat)
}

func update(r *git.Repository, quiet bool, message, tag string) error {
	if !semver.IsValid(tag) {
		return fmt.Errorf("%s is not a valid semver tag", tag)
	}

	tag = semver.Canonical(tag)
	if message == "" {
		message = "release " + tag
	}
	if !quiet {
		fmt.Printf("pushing %s: %s\n", tag, message)
		fmt.Printf("OK [y/N]? ")
		var ch string
		if n, err := fmt.Scanf("%s", &ch); n != 1 || err != nil || (ch != "Y" && ch != "y") {
			fmt.Println("Skipped....")
			return nil
		}
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	opts := &git.CreateTagOptions{Message: message}
	if opts.Tagger, err = tagger(r); err != nil {
		return err
	}

	if _, err = r.CreateTag(tag, head.Hash(), opts); err != nil {
		return err
	}

	rs := []config.RefSpec{"refs/tags/*:refs/tags/*"}
	if err = r.Push(&git.PushOptions{RemoteName: "origin", RefSpecs: rs}); err != nil {
		return err
	}

	fmt.Println("pushed tag", tag)
	return nil
}

func tagger(r *git.Repository) (*object.Signature, error) {
	c, err := r.Config()
	if err != nil {
		return nil, err
	}
	var sign object.Signature
	sign.Name = c.Raw.Section("user").Option("name")
	sign.Email = c.Raw.Section("user").Option("email")
	sign.When = time.Now()
	return &sign, nil
}

func gitDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for filepath.Dir(dir) != dir {
		f, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && f.IsDir() {
			return dir
		}
		dir = filepath.Dir(dir)
	}

	return ""
}

func help() {
	fmt.Print(`
Tag creates a semver tag in the current git repo and pushes it.

The semver versioning scheme used is something like: v5.3.2
(ie vMajor.Minor.Patch).

Tag reads the current git repository for the highest semver and
based on the command line argument passed, it increments the patch,
minor or major version.

An exact tag can also be provided for use with the very first
semver version.

If no command is provided, the current highest semver tag is
printed.
`)
	usage()
}

func usage() {
	fmt.Print(`
Usage:

tag [options] (patch|minor|major|"exact-tag")

commands:
  patch -- increments the patch version.
  major -- increments the major version.
  minor -- increments the minor version.
  tag   -- use this exact value.

options:
  -q       -- do not prompt for confirmation.
  -m <msg> -- the message to use with the tag.
  -h       -- help
`)
}
