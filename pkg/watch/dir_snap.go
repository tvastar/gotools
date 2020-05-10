package watch

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// DirSnap snapshots a dir and returns all the current files via the
// Stream. It does not watch for changes after that, returning an
// io.EOF instead. NextPath must not be called after an EOF is
// returned.
func DirSnap(root string) Stream {
	closed := make(chan error, 2) //nolint: mnd
	return &dirsnap{root, closed, nil}
}

type dirsnap struct {
	root   string
	closed chan error
	ch     chan string
}

func (d *dirsnap) NextPath(ctx context.Context) (string, error) {
	if d.ch == nil {
		d.ch = make(chan string)
		go d.walk()
	}
	select {
	case err := <-d.closed:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	case next := <-d.ch:
		return next, nil
	}
}

func (d *dirsnap) walk() {
	err := filepath.Walk(d.root, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			select {
			case <-d.closed:
				return io.EOF
			case d.ch <- path:
			}
		}
		return nil
	})

	if err == nil {
		err = io.EOF
	}
	d.closed <- err
}

func (d *dirsnap) Close() {
	d.closed <- io.EOF
}
