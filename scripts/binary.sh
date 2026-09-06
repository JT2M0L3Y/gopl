#!/bin/bash

mkdir -p dist

version="${{ steps.version.outputs.version }}"

for target in "linux amd64" "linux arm64" "darwin amd64" "darwin arm64" "windows amd64"; do
  read -r goos goarch <<< "$target"
  suffix=""

  if [[ "$goos" == "windows" ]]; then 
    suffix=".exe"; 
  fi

  name="gopl-${version}-${goos}-${goarch}"
  GOOS="$goos" GOARCH="$goarch" go build \
    -trimpath \
    -ldflags "-s -w -X gopl/internal/version.Value=${version}" \
    -o "dist/${name}${suffix}" ./cmd/gopl
done