package main

import (
	"context"
	"fmt"
	"os"

	"github.com/michimani/cfkvs/cli"
)

func main() {
	ctx := context.TODO()
	if err := cli.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
