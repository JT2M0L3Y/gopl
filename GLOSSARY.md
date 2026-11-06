# GLOSSARY.md

## Domain Terms

| Term | Definition | Source |
|------|------------|--------|
| Lexer | Converts input text into tokens | `internal/lexer/lexer.go` |
| Parser | Performs recursive descent parsing on tokens to build an AST | `internal/parser/parser.go` |
| AST | Abstract Syntax Tree representing the program structure | `internal/ast/ast.go` |
| Code Generator | Translates AST into executable code instructions | `internal/generator/generator.go` |
| Virtual Machine | Executes generated code instructions at runtime | `internal/vm/vm.go` |
| Semantic Checker | Performs type checking and semantic analysis on AST | `internal/semantic/checker.go` |
| Printer | Produces pretty-printed source code from AST | `internal/printer/print_visitor.go` |
| Pipeline | Orchestrates the full language toolchain from source to execution | `internal/pipeline/pipeline.go` |
| Token | A lexical unit with kind, lexeme, line, and column | `internal/token/token.go` |

## Architecture Decisions

| Decision | Reference |
|----------|-----------|
| Go-first port of MyPL | `README.md` |
| Preserve code generator and VM architecture | `README.md` |
| Tests beside packages as `*_test.go` files | `README.md` |
| Black-box integration tests in `tests/` directory | `README.md` |
| CLI programs under `cmd/` | `README.md` |
