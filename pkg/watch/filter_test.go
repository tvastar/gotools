package watch_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestFilter(t *testing.T) {
	s := newFixedStream([]string{"hello", "boo", "world", "hoo"})
	expected := []string{"world"}
	filter := func(s string) bool {
		return s == "world"
	}
	w := watch.Filter(filter, s)
	got, err := fetchAll(w)
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
	if err = watch.Close(w); err != nil {
		t.Error("unexpected", err)
	}
	if err = watch.Close(watch.Filter(nil, nil)); err != nil {
		t.Error("unexpected", err)
	}
}
