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
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/profile [options] <source-file>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	source, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer source.Close()

	if *cpuProfile != "" {
		file, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(file); err != nil {
			file.Close()
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer func() {
			pprof.StopCPUProfile()
			file.Close()
		}()
	}

	metrics, err := pipeline.Run(source, os.Stdin, io.Discard)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("lex=%s parse=%s semantic=%s generate=%s execute=%s\n", metrics.Lex, metrics.Parse, metrics.Semantic, metrics.Generate, metrics.Execute)

	if *memoryProfile != "" {
		file, err := os.Create(*memoryProfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(file); err != nil {
			file.Close()
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := file.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
