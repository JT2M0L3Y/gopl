# GoPL Changes

This file records major implementation decisions and the reason for each change.

## Initial direction

- Port MyPL into Go as a Go-first project instead of a direct C++ translation.
- Keep the language recognizable, but let the implementation take advantage of Go slices, maps, packages, methods, and error handling.

## Pipeline choice

- Keep the code generator and VM for now.
- Revisit that choice only after the Go port is working end to end.
- This preserves the compiler/interpreter learning value while the project is still being built.

## Design goals

- Prefer plain Go structs over deep inheritance trees.
- Prefer explicit `error` returns over exception-style control flow.
- Prefer slices and maps over custom container wrappers where possible.
- Keep runtime types and source positions clear and easy to debug.

## Planned language support

- structs
- arrays
- functions
- control flow
- builtins for printing, input, and collection helpers

## What should be documented here

Use this file for changes such as:

- AST redesign decisions
- lexer and parser behavior changes
- semantic rule changes
- runtime and VM representation choices
- CI and release pipeline changes

## What should not live here

- tiny implementation notes that belong in code comments
- day-to-day progress updates
- brainstorming that has not been turned into a decision yet

## Review rule

If we later decide the VM or code generator is no longer worth keeping, that decision should be recorded here with:

- what changed
- why it changed
- what benefit we gained
- what we gave up
