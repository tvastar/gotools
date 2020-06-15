package script

import (
	"context"
	"io"
)

type parallel []Task

func (p parallel) Start(ctx context.Context) error {
	var err error
	for _, t := range p {
		if err2 := t.Start(ctx); err2 != nil {
			err = err2
		}
	}
	return err
}

func (p parallel) Wait(ctx context.Context) error {
	var err error
	for _, t := range p {
		if err2 := t.Wait(ctx); err2 != nil {
			err = err2
		}
	}
	return err
}

func (p parallel) Stdin(r io.Reader) Task {
	result := make(parallel, len(p))
	for kk := range p {
		result[kk] = p[kk].Stdin(r)
	}
	return result
}

func (p parallel) Stdout(w io.Writer) Task {
	result := make(parallel, len(p))
	for kk := range p {
		result[kk] = p[kk].Stdout(w)
	}
	return result
}
