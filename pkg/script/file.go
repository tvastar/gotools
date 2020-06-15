package script

import (
	"context"
	"io"
	"os"
)

type file struct {
	path   string
	stdin  io.Reader
	stdout io.Writer
	f      *os.File
}

func (f *file) Start(ctx context.Context) error {
	var err error
	if f.stdin != nil {
		f.f, err = os.Create(f.path)
	} else if f.stdout != nil {
		f.f, err = os.Open(f.path)
	}
	return err
}

func (f *file) Wait(ctx context.Context) error {
	if f.f == nil {
		return nil
	}
	defer f.f.Close()

	var err error
	if f.stdin != nil {
		_, err = io.Copy(f.f, f.stdin)
	} else {
		_, err = io.Copy(f.stdout, f.f)
	}
	return err
}

func (f *file) Stdin(r io.Reader) Task {
	result := *f
	result.stdin = r
	return &result
}

func (f *file) Stdout(w io.Writer) Task {
	result := *f
	result.stdout = w
	return &result
}
