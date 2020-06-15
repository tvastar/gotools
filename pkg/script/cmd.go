package script

import (
	"context"
	"io"
	"os/exec"
	"strings"
)

type cmd struct {
	logger  Logger
	program string
	args    []string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	cmd     *exec.Cmd
}

func (c *cmd) Start(ctx context.Context) error {
	c.cmd = exec.CommandContext(ctx, c.program, c.args...)
	if c.logger != nil {
		c.logger.Println(">", strings.Join(c.cmd.Args, " "))
	}
	c.cmd.Stdin = c.stdin
	c.cmd.Stdout = c.stdout
	c.cmd.Stderr = c.stderr
	return c.cmd.Start()
}

func (c cmd) Wait(ctx context.Context) error {
	if c.cmd == nil {
		return nil
	}
	return c.cmd.Wait()
}

func (c cmd) Stdin(r io.Reader) Task {
	result := c
	result.stdin = r
	return &result
}

func (c cmd) Stdout(w io.Writer) Task {
	result := c
	result.stdout = w
	return &result
}
