// Package watch implements file watching utilities.
//
// 1. This is not concered with very large directories.
// 2. The focus is on simplicity for use in build tools.
// 3. Cross platform support would be abstracted into a singe API.
// 4. Context-aware Streams-based API.
package watch

import (
	"context"
	"os"
	"time"
)

// Stream is the main interface implemented by various watchers.
//
// It is the basis for composition (such as with Delay or Repeat or
// Filter).
type Stream interface {
	NextPath(ctx context.Context) (string, error)
}

// Dir returns a snapshot + all changes
func Dir(dir string) Stream {
	first := true
	return Dedup(LastModifiedChecksum, Repeat(func() Stream {
		if first {
			first = false
			return DirSnap(dir)
		}
		return Delay(time.Minute, DirSnap(dir))
	}))
}

// CurrentDir automatically picks the current dir but also filters
// by the glob pattern.
func CurrentDir(glob string) Stream {
	cwd, err := os.Getwd()
	if err != nil {
		return Error(err)
	}
	first := true
	return Dedup(LastModifiedChecksum, Repeat(func() Stream {
		filtered := Filter(Glob(glob), DirSnap(cwd))
		if first {
			first = false
			return filtered
		}
		return Delay(time.Minute, filtered)
	}))
}
