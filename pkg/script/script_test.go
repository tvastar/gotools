package script_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tvastar/gotools/pkg/script"
)

func ExampleSequence() {
	logger := log.New(os.Stdout, "", 0)
	spec := script.Sequence(
		script.CmdWithLog(logger, "echo", "hello"),
		script.CmdWithLog(logger, "echo", "world"),
	)
	err := script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	// Output:
	// > echo hello
	// hello
	// > echo world
	// world
}

func ExampleFile_redirect() {
	spec := script.Sequence(
		// echo hello > /tmp/hello.txt
		script.Pipe(
			script.Cmd("echo", "hello"),
			script.File("/tmp/hello.txt"),
		),

		// cat < /tmp/hello.txt
		script.Pipe(
			script.File("/tmp/hello.txt"),
			script.Cmd("cat"),
		),
		// script.Cmd("rm", "/tmp/hello.txt"),
	)
	err := script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	// Output: hello
}

func ExampleIf() {
	spec := script.If(
		script.Cmd("false"),
		script.Cmd("echo", "wrongly succeeded"),
		script.Cmd("echo", "rightly failed"),
	)
	err := script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	spec = script.If(
		script.IgnoreError(script.Cmd("false")),
		script.Cmd("echo", "rightly succeeded"),
		script.Cmd("echo", "wrongly failed"),
	)
	err = script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	// Output:
	// rightly failed
	// rightly succeeded
}

func ExampleOr() {
	spec := script.Or(
		script.Cmd("false"),
		script.Cmd("false"),
		script.Cmd("echo", "hello"),
	)
	err := script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	// Output: hello
}

func ExamplePipe_if() {
	spec := script.If(
		script.Pipe(
			script.Cmd("echo", "hello"),
			script.Cmd("sed", "s/hello/world/"),
		),
		script.Cmd("echo", "rightly succeeded"),
		script.Cmd("echo", "wrongly failed"),
	)
	err := script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	// Output:
	// world
	// rightly succeeded
}

func ExampleFunc() {
	spec := script.Pipe(
		script.Cmd("echo", "hello"),
		script.Func(func(ctx context.Context, r io.Reader, w io.Writer) error {
			_, err := io.Copy(w, r)
			return err
		}),
		script.Cmd("cat"),
	)
	err := script.Run(context.Background(), spec)
	if err != nil {
		fmt.Println("error", err)
	}

	// Output: hello
}
