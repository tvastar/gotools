package watch

import "context"

// Error implements a stream that always returns errors.
func Error(err error) Stream {
	return errStream{err}
}

type errStream struct {
	err error
}

func (e errStream) NextPath(ctx context.Context) (path string, err error) {
	return "", e.err
}
