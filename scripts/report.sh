#!/bin/bash
set -euo pipefail

{
  echo "# GoPL v${version}"
  echo
  echo "## Performance sign-off"
  echo
  echo "These results were collected from merge commit ${merge_sha} before publication. Review the benchmark and profile artifacts before approving the release."
  echo
  echo '```text'
  cat artifacts/benchmark.txt
  echo '```'
  echo
  echo "Profiles and raw benchmark output are attached to this release."
} > artifacts/release-notes.md