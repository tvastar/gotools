package script

import (
	"context"
	"io"
)

type fn struct {
	f       func(ctx context.Context) error
	started bool
}

func (f *fn) Start(ctx context.Context) error {
	f.started = true
	return nil
}

func (f fn) Wait(ctx context.Context) error {
	if f.started {
		return f.f(ctx)
	}
	return nil
}

func (f fn) Stdin(r io.Reader) Task {
	panic("cannot pipe into to a function")
}

func (f fn) Stdout(w io.Writer) Task {
	panic("cannot pipe out of a function")
}
