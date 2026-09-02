package main

import (
	"flag"
	"fmt"
	"os"

	"gopl/internal/pipeline"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: gopl <source-file>")
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	src, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer src.Close()
	if _, err := pipeline.Run(src, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
