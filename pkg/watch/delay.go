package watch

import (
	"context"
	"time"
)

// Delay waits for the specified time before allowing the first call
// to proceed.  This is useful when combined with Repeat to setup a
// polling interval.
func Delay(duration time.Duration, s Stream) Stream {
	return &delay{time.NewTimer(duration), s}
}

type delay struct {
	*time.Timer
	s Stream
}

func (d *delay) NextPath(ctx context.Context) (string, error) {
	if d.Timer != nil {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-d.Timer.C:
		}
		d.Timer.Stop()
		d.Timer = nil
	}
	return d.s.NextPath(ctx)
}

func (d *delay) Close() error {
	if d.Timer != nil {
		d.Timer.Stop()
		d.Timer = nil
	}
	return Close(d.s)
}
