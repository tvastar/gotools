package watch

import "io"

// Close closes a stream by checking if the stream implements io.Closer.
func Close(s Stream) error {
	if closer, ok := s.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
