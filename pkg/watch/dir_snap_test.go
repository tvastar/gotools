package watch_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestDirSnap(t *testing.T) {
	got, err := fetchAll(watch.DirSnap("testdata"))
	expected := []string{"testdata", "one.txt", "two.txt"}
	if len(got) != len(expected) || err != io.EOF {
		t.Error("unexpedted", got)
	}

	for idx, g := range got {
		if !strings.HasSuffix(g, expected[idx]) {
			t.Error("no suffix", g, expected[idx])
		}
	}
}

func TestDirSnapPartialClose(t *testing.T) {
	w := watch.DirSnap("testdata")
	if _, err := w.NextPath(context.Background()); err != nil {
		t.Fatal("unexpected", err)
	}
	if err := watch.Close(w); err != nil {
		t.Fatal("close failed", err)
	}
}
