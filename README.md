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

Tests live beside the package they exercise as `*_test.go` files. The `tests/` directory is reserved for black-box integration tests and source fixtures. Run the complete suite with `go test ./...`.

The command-line programs live under `cmd/`: run the interpreter with `go run ./cmd/gopl <source-file>` and the profiler with `go run ./cmd/profile [options] <source-file>`.

## Quality checks

Local validation mirrors CI:

```text
gofmt -l .
go vet ./...
go test ./...
go build ./...
```

CI also runs `golangci-lint` on every push and pull request. Scheduled and manually triggered performance workflows run repeated benchmarks and upload CPU and heap profiles. Profiles can be generated locally with:

```text
go run ./cmd/profile -cpuprofile gopl.cpu.pprof -memprofile gopl.mem.pprof tests/fixtures/profile.gopl
```
