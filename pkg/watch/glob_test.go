package watch_test

import (
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestGlob(t *testing.T) {
	glob := watch.Glob("**/boo")
	if !glob("/tmp/foo/boo") || glob("/tmp/boorish") {
		t.Error("unexpected glob behaviors")
	}
}
