package main

import (
	"flag"
	"fmt"
	"os"

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
	parseMode := flag.Bool("parse", false, "placeholder for the parser stage")
	printMode := flag.Bool("print", false, "placeholder for the pretty-printer stage")
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
	case modeParse, modePrint:
		fmt.Fprintln(os.Stderr, "that stage is not implemented yet")
		os.Exit(1)
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
