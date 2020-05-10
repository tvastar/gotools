// +build darwin

package watch_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestFSEvents(t *testing.T) {
	dir, err := ioutil.TempDir("", "fsevents_test")
	if err != nil {
		t.Fatal("create temp dir", err)
	}
	defer os.RemoveAll(dir)

	if dir, err = filepath.EvalSymlinks(dir); err != nil {
		t.Fatal("evalsymlinks", err)
	}

	w := watch.DirFSEvents(dir)
	defer watch.Close(w)

	if s, err := w.NextPath(context.Background()); s != dir || err != nil {
		t.Error("unexpected path", s, dir)
	}

	fname := filepath.Join(dir, "somefile.txt")
	if err := ioutil.WriteFile(fname, []byte("something"), 0666); err != nil {
		t.Error("writefile", err)
	}

	for {
		s, err := w.NextPath(context.Background())
		switch {
		case s == dir && err == nil:
		case s == fname && err == nil:
			return
		default:
			t.Fatal("unexpected result", s, err)
		}
	}
}
