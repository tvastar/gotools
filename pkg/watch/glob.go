package watch

import (
	"github.com/bmatcuk/doublestar"
)

// Glob tests if strings match the specified pattern.
func Glob(pattern string) (allow func(path string) bool) {
	return func(path string) bool {
		ok, err := doublestar.PathMatch(pattern, path)
		return ok && err == nil
	}
}
