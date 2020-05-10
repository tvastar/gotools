// Command watch prints file names as they change.
//
// Usage:
//
//   watch [optional_global_pattern]
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tvastar/gotools/pkg/watch"
)

func main() {
	ctx := context.Background()
	glob := "**"
	if len(os.Args) > 1 {
		glob = os.Args[1]
	}

	w := watch.CurrentDir(glob)

	for {
		p, err := w.NextPath(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error", err)
			os.Exit(1)
		}

		fmt.Println(p)
	}
}
