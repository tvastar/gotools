package watch_test

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/tvastar/gotools/pkg/watch"
)

func TestDelayNoDelay(t *testing.T) {
	s := newFixedStream([]string{"hello", "boo", "world", "hoo"})
	expected := []string{"hello", "boo", "world", "hoo"}
	got, err := fetchAll(watch.Delay(0, s))
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
}

func TestDelaySmallDelay(t *testing.T) {
	s := newFixedStream([]string{"hello", "boo", "world", "hoo"})
	expected := []string{"hello", "boo", "world", "hoo"}
	got, err := fetchAll(watch.Delay(time.Millisecond, s))
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
}

func TestDelayCancelation(t *testing.T) {
	s := newFixedStream([]string{"hello", "boo", "world", "hoo"})
	expected := []string{"hello", "boo", "world", "hoo"}
	delayed := watch.Delay(time.Millisecond*5, s)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	if s, err := delayed.NextPath(ctx); err == nil || err != ctx.Err() {
		t.Fatal("unexpected delay", s, err)
	}
	got, err := fetchAll(delayed)
	if !reflect.DeepEqual(expected, got) || err != io.EOF {
		t.Error("unexpected", got, err)
	}
}
