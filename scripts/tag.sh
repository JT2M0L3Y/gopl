#!/bin/bash
set -euo pipefail

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git tag -a "v${version}" "${merge_sha}" -m "Release v${version}"
git push origin "v${version}"