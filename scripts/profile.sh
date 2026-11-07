#!/bin/bash
set -euo pipefail

mkdir -p artifacts

go run ./cmd/profile \
  -cpuprofile artifacts/gopl-${version}.cpu.pprof \
  -memprofile artifacts/gopl-${version}.mem.pprof \
  tests/fixtures/profile.gopl