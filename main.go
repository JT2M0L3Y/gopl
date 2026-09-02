package main

import (
	"flag"
	"fmt"
	"os"

	"gopl/internal/generator"
	"gopl/internal/lexer"
	"gopl/internal/parser"
	"gopl/internal/semantic"
	"gopl/internal/vm"
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
	runProgram(src)
}

func runProgram(src *os.File) {
	l := lexer.New(src)
	p := parser.New(l)
	prog, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	checker := semantic.NewSemanticChecker()
	if err := checker.Check(prog); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	runtime := vm.New()
	code := generator.New(runtime)
	if err := code.Generate(prog); err != nil {
		fmt.Fprintf(os.Stderr, "code generation error: %v\n", err)
		os.Exit(1)
	}
	if err := runtime.Execute(false); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
