package watch

import "context"

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
	return Close(f.s)
}
