#!/bin/bash

go run ./cmd/profile \
  -cpuprofile artifacts/gopl-${{ steps.version.outputs.version }}.cpu.pprof \
  -memprofile artifacts/gopl-${{ steps.version.outputs.version }}.mem.pprof tests/fixtures/profile.gopl