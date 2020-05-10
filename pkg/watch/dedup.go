package watch

import "os"

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
	checksums := map[string]interface{}{}
	allow := func(path string) bool {
		current := checksum(path)
		old, ok := checksums[path]
		if ok && old == current {
			return false
		}

		if current == nil {
			delete(checksums, path)
		} else {
			checksums[path] = current
		}
		return true
	}
	return Filter(allow, s)
}
