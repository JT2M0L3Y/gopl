package main

import (
	"flag"
	"fmt"
	"os"

	"gopl/internal/ast"
	"gopl/internal/lexer"
	"gopl/internal/token"
)

type mode string

const (
	modeLex   mode = "lex"
	modeParse mode = "parse"
	modePrint mode = "print"
)

func main() {
	lexMode := flag.Bool("lex", false, "run the lexer and print tokens")
	parseMode := flag.Bool("parse", false, "run the parser and print parse info")
	printMode := flag.Bool("print", false, "run the parser and pretty-print the AST")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: gopl <source-file>  # default mode is lex")
		flag.PrintDefaults()
	}
	flag.Parse()

	m, err := selectedMode(*lexMode, *parseMode, *printMode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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

	switch m {
	case modeLex:
		runLexer(src)
	case modeParse:
		runParser(src)
	case modePrint:
		runPrinter(src)
	}
}

func selectedMode(lexMode, parseMode, printMode bool) (mode, error) {
	count := 0
	var m mode
	if lexMode {
		m = modeLex
		count++
	}
	if parseMode {
		m = modeParse
		count++
	}
	if printMode {
		m = modePrint
		count++
	}
	if count == 0 {
		return modeLex, nil
	}
	if count > 1 {
		return "", fmt.Errorf("choose only one of --lex, --parse, or --print")
	}
	return m, nil
}

func runLexer(src *os.File) {
	l := lexer.New(src)
	for {
		tok := l.Next()
		fmt.Println(tok.String())
		if tok.Kind == token.EOF {
			break
		}
	}
}

func runParser(src *os.File) {
	l := lexer.New(src)
	p := ast.New(l)
	prog, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Successfully parsed program:\n")
	fmt.Printf("  Structs: %d\n", len(prog.StructDefs))
	fmt.Printf("  Functions: %d\n", len(prog.FunDefs))
}

func runPrinter(src *os.File) {
	l := lexer.New(src)
	p := ast.New(l)
	prog, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	pv := ast.NewPrintVisitor(os.Stdout)
	if err := prog.Accept(pv); err != nil {
		fmt.Fprintf(os.Stderr, "visitor error: %v\n", err)
		os.Exit(1)
	}
}
