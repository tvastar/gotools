package watch_test

import (
	"context"
	"io"

	"github.com/tvastar/gotools/pkg/watch"
)

type fakestream struct {
	nextPath func() (string, error)
	close    func() error
}

func (f *fakestream) NextPath(ctx context.Context) (string, error) {
	return f.nextPath()
}

func (f *fakestream) Close() error {
	return f.close()
}

func newFixedStream(inputs []string) watch.Stream {
	return &fakestream{
		nextPath: func() (string, error) {
			if len(inputs) == 0 {
				return "", io.EOF
			}
			next := inputs[0]
			inputs = inputs[1:]
			return next, nil
		},
		close: func() error {
			return nil
		},
	}
}

func fetchAll(s watch.Stream) ([]string, error) {
	paths := []string{}
	for {
		p, err := s.NextPath(context.Background())
		if err != nil {
			return paths, err
		}
		paths = append(paths, p)
	}
}
