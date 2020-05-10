package watch_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestRepeat(t *testing.T) {
	someErr := errors.New("some error")

	streams := []watch.Stream{
		newFixedStream([]string{"hello", "world"}),
		newFixedStream([]string{"boo", "hoo"}),
		watch.Error(someErr),
	}
	create := func() watch.Stream {
		next := streams[0]
		streams = streams[1:]
		return next
	}

	expected := []string{"hello", "world", "boo", "hoo"}
	w := watch.Repeat(create)
	got, err := fetchAll(w)
	if !reflect.DeepEqual(expected, got) || err != someErr {
		t.Error("unexpected", got, err)
	}
	if err = watch.Close(w); err != nil {
		t.Error("unexpected", err)
	}
}
