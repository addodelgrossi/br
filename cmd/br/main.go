package main

import (
	"fmt"
	"os"

	"github.com/addodelgrossi/br/internal/cli"
)

func main() {
	root := cli.NewRootCommand(os.Stdout)
	root.SetArgs(cli.PrepareArgs(os.Args[1:]))

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
