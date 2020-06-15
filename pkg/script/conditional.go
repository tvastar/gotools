package script

import (
	"context"
	"io"
)

type conditional struct {
	cond, success, failure Task
	started                bool
	err                    error
}

func (c *conditional) Start(ctx context.Context) error {
	c.err = c.cond.Start(ctx)
	c.started = true
	return c.err
}

func (c *conditional) Wait(ctx context.Context) error {
	if !c.started {
		return nil
	}

	result := c.success
	if c.err == nil {
		c.err = c.cond.Wait(ctx)
	}
	if c.err != nil {
		result = c.failure
	}

	if result == nil {
		return c.err
	}
	return Run(ctx, result)
}

func (c *conditional) Stdin(r io.Reader) Task {
	withStdin := func(t Task) Task {
		if t != nil {
			return t.Stdin(r)
		}
		return t
	}

	return &conditional{
		cond:    c.cond.Stdin(r),
		success: withStdin(c.success),
		failure: withStdin(c.failure),
	}
}

func (c *conditional) Stdout(w io.Writer) Task {
	withStdout := func(t Task) Task {
		if t != nil {
			return t.Stdout(w)
		}
		return t
	}

	return &conditional{
		cond:    c.cond.Stdout(w),
		success: withStdout(c.success),
		failure: withStdout(c.failure),
	}
}
