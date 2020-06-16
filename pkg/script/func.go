package script

import (
	"context"
	"io"
)

type fn struct {
	f      func(ctx context.Context, r io.Reader, w io.Writer) error
	ch     chan error
	stdin  io.Reader
	stdout io.Writer
}

func (f *fn) Start(ctx context.Context) error {
	f.ch = make(chan error, 1)
	go func() {
		var err error
		defer func() {
			f.ch <- err
		}()
		err = f.f(ctx, f.stdin, f.stdout)
	}()
	return nil
}

func (f fn) Wait(ctx context.Context) error {
	if f.ch != nil {
		return <-f.ch
	}
	return nil
}

func (f fn) Stdin(r io.Reader) Task {
	result := f
	result.stdin = r
	return &result
}

func (f fn) Stdout(w io.Writer) Task {
	result := f
	result.stdout = w
	return &result
}
