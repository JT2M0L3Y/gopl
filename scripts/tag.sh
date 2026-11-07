#!/bin/bash
set -euo pipefail

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git tag -a "$version" "${merge_sha}" -m "Release $version"
git push origin "v${version}"