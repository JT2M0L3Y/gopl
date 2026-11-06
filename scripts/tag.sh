#!/bin/bash

version="v${{ needs.performance.outputs.version }}"

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git tag -a "$version" "${{ needs.performance.outputs.merge_sha }}" -m "Release $version"

git push origin "$version"