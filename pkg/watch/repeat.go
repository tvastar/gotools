package watch

import (
	"context"
	"io"
)

// Repeat calls the create function repeated, iterating through each
// one until an io.EOF is returned.
func Repeat(create func() Stream) Stream {
	return &repeat{create, nil}
}

type repeat struct {
	create func() Stream
	Stream
}

func (p *repeat) NextPath(ctx context.Context) (string, error) {
	for {
		if p.Stream == nil {
			p.Stream = p.create()
		}
		if s, err := p.Stream.NextPath(ctx); err != io.EOF {
			return s, err
		}
		p.Stream = nil
	}
}

func (p *repeat) Close() error {
	return Close(p.Stream)
}
