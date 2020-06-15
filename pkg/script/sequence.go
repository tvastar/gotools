package script

import (
	"context"
	"io"
)

type seq []Task

func (s seq) Start(ctx context.Context) error {
	if len(s) == 0 {
		return nil
	}
	return s[0].Start(ctx)
}

func (s seq) Wait(ctx context.Context) error {
	for kk := range s {
		if kk > 0 {
			if err := s[kk].Start(ctx); err != nil {
				return err
			}
		}
		if err := s[kk].Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s seq) Stdin(r io.Reader) Task {
	result := make(seq, len(s))
	for kk := range s {
		result[kk] = s[kk].Stdin(r)
	}
	return result
}

func (s seq) Stdout(w io.Writer) Task {
	result := make(seq, len(s))
	for kk := range s {
		result[kk] = s[kk].Stdout(w)
	}
	return result
}
