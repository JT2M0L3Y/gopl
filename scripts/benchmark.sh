#!/bin/bash

mkdir -p artifacts

go test ./internal/pipeline -run '^$' -bench BenchmarkPipelineStages -benchmem -count=1 | tee -a artifacts/benchmark.txt