// Package script makes it easy to write declarative scripts in go.
//
// The core abstraction is a Task. A shell command is an example of a
// task and this can be created and executed via:
//
//     task := script.Cmd("echo", "hello")
//     err := script.Run(context.Background(), task)
//
// Tasks can be composed in sequence:
//
//     task := script.Sequence(
//         script.Cmd("echo", "hello"),
//         script.Cmd("echo", "world"),
//     )
//     err := script.Run(context.Background(), task)
//
// This type of composition allows for expressions without having to
// constantly check for errors.  A sequential task stops with an error
// when one of its items fails.
//
// The other types of control flow structures are "Parallel", "Or" and
// "If".
//
// Tasks can also be "Piped" in a chain via `Pipe`:
//
//      task := script.Pipe(
//          script.Cmd("echo", "hello"),
//          script.Cmd("sed", "s/hello/world"),
//          script.Cmd("grep", "world"),
//      )
//      err := script.Run(context.Background(), task)
//
// These can all be nested to express the intent.
//
// Tasks can also be piped to a file (or piped from a file) by just
// using `File(path)` in a pipe.
//
// Non shell tasks can be mingled within via the `Func(...)` method.
package script

import (
	"context"
	"io"
	"os"
)

// Task defines a task that can be run.
type Task interface {
	// Start starts a task asynchronously.
	Start(ctx context.Context) error

	// Wait waits for the task to finish.
	// Wait can be called before the task is started. In this
	// case, it should do nothing and return nil.
	Wait(ctx context.Context) error

	// Stdin redirects input.
	// This is immutable returning a new task.
	Stdin(r io.Reader) Task

	// Stdout redirects output.
	// This is immutable returning a new task.
	Stdout(w io.Writer) Task
}

// Run runs a task.
func Run(ctx context.Context, t Task) error {
	err1 := t.Start(ctx)
	err2 := t.Wait(ctx)
	if err1 != nil {
		return err1
	}
	return err2
}

// Sequence chains a sequence of tasks together. If any task fails,
// the sequence is aborted.
func Sequence(tasks ...Task) Task {
	return seq(tasks)
}

// Parallel runs all tasks in parallel. If any tasks fail, it returns
// one of the errors.
func Parallel(tasks ...Task) Task {
	return parallel(tasks)
}

// Pipe runs all tasks in a "pipe" chaining their input and outputs.
func Pipe(tasks ...Task) Task {
	writers := make([]*os.File, len(tasks)-1)
	readers := make([]*os.File, len(tasks)-1)
	return pipe{tasks, writers, readers}
}

// File runs a task which allows either input to be piped to it or for
// the file to be piped into another task.
func File(path string) Task {
	return &file{path: path}
}

// Or runs a sequence of commands until one succeeds.
func Or(tasks ...Task) Task {
	var result Task
	for kk := len(tasks) - 1; kk >= 0; kk-- {
		if kk == len(tasks)-1 {
			result = tasks[kk]
		} else {
			result = If(tasks[kk], nil, result)
		}
	}
	return result
}

// IgnoreError runs a task and ignores any errors.
func IgnoreError(t Task) Task {
	// use an empty sequence to ignore errors.
	return If(t, nil, Sequence())
}

// If runs the "then" task if the "condition" succeeds. It runs the
// else task if the "condition" fails.
// It is legal for thenTask and elseTask to be nil.  If either of them
// are nil, the result of the condition is propagated.
func If(condition, thenTask, elseTask Task) Task {
	return &conditional{condition, thenTask, elseTask, false, nil}
}

// Cmd runs a program with the provided args.
func Cmd(program string, args ...string) Task {
	return &cmd{nil, program, args, os.Stdin, os.Stdout, os.Stderr, nil}
}

type Logger interface {
	Println(v ...interface{})
}

// CmdWithLog runs a program with the provided args and also logs output.
func CmdWithLog(logger Logger, program string, args ...string) Task {
	return &cmd{logger, program, args, os.Stdin, os.Stdout, os.Stderr, nil}
}

// Func runs a separate task function.  Note that this function cannot
// participate in pipes and such.
func Func(f func(ctx context.Context) error) Task {
	return &fn{f, false}
}
