package watch

import (
	"context"
	"os"
)

// LastModifiedChecksum uses the last modified time as the checksum
// for a path.
func LastModifiedChecksum(path string) interface{} {
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return fi.ModTime()
}

// Dedup drops duplicates in the stream.
//
// Duplicates are identified by comparing the results of the checksum
// against previous results, if any.
//
// If the checksum returns nil, the last checksum is uncached.
func Dedup(checksum func(string) interface{}, s Stream) Stream {
	return &dedup{
		checksum:  checksum,
		checksums: map[string]interface{}{},
		s:         s,
	}
}

type dedup struct {
	checksum  func(string) interface{}
	checksums map[string]interface{}
	s         Stream
}

func (d dedup) NextPath(ctx context.Context) (string, error) {
	for {
		path, err := d.s.NextPath(ctx)
		if err != nil {
			return "", err
		}
		current := d.checksum(path)
		if old, ok := d.checksums[path]; !ok || old != current {
			if current == nil {
				delete(d.checksums, path)
			} else {
				d.checksums[path] = current
			}
			return path, nil
		}
	}
}

func (d dedup) Close() error {
	return Close(d.s)
}
