package watch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestError(t *testing.T) {
	someErr := errors.New("some error")
	s := watch.Error(someErr)

	if _, err := s.NextPath(context.Background()); err != someErr {
		t.Error("unexpected result", err)
	}

	if err := watch.Close(s); err != nil {
		t.Error("unexpected error", err)
	}
}
