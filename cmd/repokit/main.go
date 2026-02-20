package main

import (
	"fmt"
	"os"
)

// Version is set via ldflags at build time.
var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// TODO: implement cobra CLI
	fmt.Println("repokit", version)
	return nil
}
