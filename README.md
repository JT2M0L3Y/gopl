# GoPL

GoPL is a Go-first port of MyPL, a small programming language interpreter originally written in C++.

The goal is not a line-for-line translation. The goal is to learn Go while building a complete language toolchain:

- lexer and tokenizer
- parser and AST
- pretty printer
- semantic analysis
- symbol and variable tables
- code generator
- virtual machine

We will keep the code generator and VM for now so the project preserves the full compiler/interpreter pipeline. Once the Go implementation is complete, we can decide whether that architecture still earns its complexity or whether the runtime should be simplified.

## Roadmap

The working roadmap lives in [`STEPS.md`](./STEPS.md).

Implementation decisions and architecture notes live in [`CHANGES.md`](./CHANGES.md).

## Working style

The port will move in layers:

1. copy the existing ideas into Go
2. reshape the code to fit Go idioms
3. add or remove pieces only when Go benefits from the change
4. document each stage so the project doubles as a learning resource

## Current focus

The current implementation walks from the entrypoint through tokenization, parsing, AST construction, semantic checks, code generation, and VM execution. Additional language features will be added deliberately as separate stages.
