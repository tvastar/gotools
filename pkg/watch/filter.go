package watch

import (
	"context"
	"io"
)

// Filter only returns paths matching the allow filter.
func Filter(allow func(path string) bool, s Stream) Stream {
	return filter{allow, s}
}

type filter struct {
	allow func(path string) bool
	s     Stream
}

func (f filter) NextPath(ctx context.Context) (string, error) {
	for {
		p, err := f.s.NextPath(ctx)
		if err == nil && !f.allow(p) {
			continue
		}
		return p, err
	}
}

func (f filter) Close() error {
	if closer, ok := f.s.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
