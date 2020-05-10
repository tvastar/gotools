package handles_test

import (
	"testing"

	"github.com/tvastar/gotools/pkg/handles"
)

func TestHandles(t *testing.T) {
	var htable handles.Table

	// empty get
	if _, ok := htable.Get(0); ok {
		t.Fatal("empty get succeeded")
	}

	// add
	h := htable.Add("hello")
	if v, ok := htable.Get(h); !ok || v != "hello" {
		t.Fatal("Add failed", h, v, ok)
	}

	// remove
	if !htable.Delete(h) || htable.Size() != 0 {
		t.Fatal("delete failed", htable.Size())
	}

	// fetch deleted
	if v, ok := htable.Get(h); ok {
		t.Fatal("unexpected fetch of deleted", v, ok)
	}
}
