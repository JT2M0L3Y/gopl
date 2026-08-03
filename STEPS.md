# GoPL Steps

This document is the live roadmap for the MyPL to GoPL port.

The main goals are:

- learn Go through a real compiler/interpreter project
- preserve the useful parts of the MyPL language
- redesign the implementation to fit Go idioms
- complete dictionary and map support
- keep the code generator and VM for now, then revisit the architecture once the Go version is working end to end

## Stage 1: Project shape

- define the Go module layout
- replace the placeholder `gopl.go` entrypoint with a real CLI
- decide the package boundaries for lexer, parser, AST, semantic analysis, code generation, VM, and runtime support
- document the migration path in this repository

## Stage 2: Tokens and lexing

- port token definitions
- implement source location tracking
- build a Go lexer that produces tokens from input text
- make lexer errors clear and idiomatic for Go

## Stage 3: Parsing and AST

- port the parser
- define AST nodes as Go structs and slices
- keep the tree easy to inspect and print
- make the parser own syntax errors and recovery decisions where practical

## Stage 4: Pretty printing

- port the pretty printer
- use it as a sanity check for the parser and AST
- keep formatting consistent and readable

## Stage 5: Semantic analysis

- port the symbol table
- port the variable table
- port type and scope checks
- add or refine dictionary and map typing rules

## Stage 6: Code generation

- port the bytecode or instruction model
- generate machine instructions from the AST
- keep the generator understandable before making it clever
- confirm the generator works with the Go data structures

## Stage 7: Virtual machine

- port the VM runtime
- implement stack, call frame, heap, and object handling
- support structs, arrays, and dictionaries
- confirm runtime behavior matches the language rules

## Stage 8: Builtins and language completeness

- finish dictionary and map methods
- verify string and array helpers
- implement remaining builtins and edge cases
- decide whether any original features should be simplified for Go

## Stage 9: Tests

- add unit tests around the lexer
- add parser and AST tests
- add semantic tests
- add VM and end-to-end execution tests

## Stage 10: GitHub Actions

- build on every push and pull request
- run the test suite automatically
- publish release artifacts when a tagged release is created
- keep the workflow simple enough to understand and maintain

## Stage 11: Architecture review

- evaluate whether the code generator and VM still feel worthwhile
- compare the current architecture against a simpler direct interpreter
- decide whether any layer should be removed, collapsed, or kept

## Learning Notes

As we work, each stage should explain:

- what the original C++ version did
- how the Go version differs
- why the Go choice is better here
- what tradeoffs we accepted

That way the repo acts as both an implementation and a Go learning guide.
