package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/pprof"

	"gopl/internal/pipeline"
)

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write a CPU profile to this file")
	memoryProfile := flag.String("memprofile", "", "write a heap profile to this file")
	flag.Usage = func() {
		writeStderr("usage: go run ./cmd/profile [options] <source-file>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	source, err := os.Open(flag.Arg(0))
	if err != nil {
		writeStderr(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := source.Close(); err != nil {
			writeStderr(err.Error())
		}
	}()

	if *cpuProfile != "" {
		file, err := os.Create(*cpuProfile)
		if err != nil {
			writeStderr(err.Error())
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(file); err != nil {
			if closeErr := file.Close(); closeErr != nil {
				writeStderr(closeErr.Error())
			}
			writeStderr(err.Error())
			os.Exit(1)
		}
		defer func() {
			pprof.StopCPUProfile()
			if err := file.Close(); err != nil {
				writeStderr(err.Error())
			}
		}()
	}

	metrics, err := pipeline.Run(source, os.Stdin, io.Discard)
	if err != nil {
		writeStderr(err.Error())
		os.Exit(1)
	}
	if _, err := fmt.Printf("lex=%s parse=%s semantic=%s generate=%s execute=%s\n", metrics.Lex, metrics.Parse, metrics.Semantic, metrics.Generate, metrics.Execute); err != nil {
		os.Exit(1)
	}

	if *memoryProfile != "" {
		file, err := os.Create(*memoryProfile)
		if err != nil {
			writeStderr(err.Error())
			os.Exit(1)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(file); err != nil {
			if closeErr := file.Close(); closeErr != nil {
				writeStderr(closeErr.Error())
			}
			writeStderr(err.Error())
			os.Exit(1)
		}
		if err := file.Close(); err != nil {
			writeStderr(err.Error())
			os.Exit(1)
		}
	}
}

func writeStderr(message string) {
	if _, err := fmt.Fprintln(os.Stderr, message); err != nil {
		os.Exit(1)
	}
}
