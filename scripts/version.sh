#!/bin/bash

latest=$(git tag --list 'v[0-9]*' --sort=-version:refname | head -n 1)

if [[ -z "$latest" ]]; then
  version="0.1.0"
else
  current="${latest#v}"
  IFS=. read -r major minor patch <<< "$current"
  version="${major}.${minor}.$((patch + 1))"
fi

if git rev-parse "v${version}" >/dev/null 2>&1; then
  echo "version v${version} already exists" >&2
  exit 1
fi

echo "version=${version}" >> "$GITHUB_OUTPUT"
echo "Next release: v${version}"
