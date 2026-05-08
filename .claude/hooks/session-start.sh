#!/bin/bash
set -euo pipefail

echo '{"async": true, "asyncTimeout": 300000}'

if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  echo "session-start: skipping setup (CLAUDE_CODE_REMOTE not set)" >&2
  exit 0
fi

"$(dirname "$0")/../../bin/setup"
