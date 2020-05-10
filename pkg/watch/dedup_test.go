package watch_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestLastModifiedChecksum(t *testing.T) {
	if v := watch.LastModifiedChecksum("goop"); v != nil {
		t.Error("unexpected last modified checksum", v)
	}
	v1 := watch.LastModifiedChecksum("dedup_test.go")
	v2 := watch.LastModifiedChecksum("dedup_test.go")
	if v1 == nil || v1 != v2 {
		t.Error("unexpected last modified checksum", v1, v2)
	}
}

func TestDedupUniques(t *testing.T) {
	counter := 0
	checksum := func(s string) interface{} {
		counter++
		return counter
	}

	s := newFixedStream([]string{"hello", "boo", "world", "hoo"})
	expected := []string{"hello", "boo", "world", "hoo"}
	got, err := fetchAll(watch.Dedup(checksum, s))
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
}

func TestDedupNoCache(t *testing.T) {
	checksum := func(s string) interface{} {
		return nil
	}

	s := newFixedStream([]string{"hello", "boo", "world", "hoo"})
	expected := []string{"hello", "boo", "world", "hoo"}
	got, err := fetchAll(watch.Dedup(checksum, s))
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
}

func TestDedupWithDupes(t *testing.T) {
	checksum := func(s string) interface{} {
		return s
	}

	s := newFixedStream([]string{"hello", "boo", "hello", "hoo"})
	expected := []string{"hello", "boo", "hoo"}
	got, err := fetchAll(watch.Dedup(checksum, s))
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
}
