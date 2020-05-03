package dochtml

import (
	"errors"
	"flag"
	"regexp"
	"strings"
)

// type assertion
var _ flag.Value = &FileLinker{}

// FileLinker linkifies source files.
type FileLinker struct {
	matches      []*regexp.Regexp
	replacements []string
}

// String is not yet implemented.
func (f *FileLinker) String() string {
	return "NYI"
}

// Set parses a regexp=replacemennt pattern and appends to the list.
func (f *FileLinker) Set(v string) error {
	idx := strings.Index(v, "=")
	if idx < 0 {
		return errors.New("src must be regexp=replacement")
	}
	re, err := regexp.Compile(v[:idx])
	if err != nil {
		return err
	}
	f.matches = append(f.matches, re)
	f.replacements = append(f.replacements, v[idx+1:])
	return nil
}

// URL converts a file path to an url.
func (f *FileLinker) URL(file string) (string, bool) {
	for idx := range f.matches {
		r := f.matches[idx].ReplaceAllString(file, f.replacements[idx])
		if r != file {
			return r, true
		}
	}
	return "", false
}
