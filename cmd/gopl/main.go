package main

import (
	"flag"
	"fmt"
	"os"

	"gopl/internal/pipeline"
	"gopl/internal/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "print the GoPL version")
	flag.Usage = func() {
		writeStderr("usage: gopl <source-file>")
	}
	flag.Parse()
	if *versionFlag {
		if _, err := fmt.Fprintln(os.Stdout, version.String()); err != nil {
			writeStderr(err.Error())
			os.Exit(1)
		}
		return
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	src, err := os.Open(flag.Arg(0))
	if err != nil {
		writeStderr(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := src.Close(); err != nil {
			writeStderr(err.Error())
		}
	}()
	if _, err := pipeline.Run(src, os.Stdin, os.Stdout); err != nil {
		writeStderr(err.Error())
		os.Exit(1)
	}
}

func writeStderr(message string) {
	if _, err := fmt.Fprintln(os.Stderr, message); err != nil {
		os.Exit(1)
	}
}
