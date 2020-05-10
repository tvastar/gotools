package watch_test

import (
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestClose(t *testing.T) {
	var closed bool

	f := fakestream{close: func() error {
		closed = true
		return nil
	}}
	if err := watch.Close(&f); err != nil || !closed {
		t.Fatal("unexpected error", err, closed)
	}
}
