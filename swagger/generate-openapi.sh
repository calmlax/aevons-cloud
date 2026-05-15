#!/usr/bin/env bash

set -euo pipefail

SERVICE_ROOT="${1:-.}"
ENTRY_FILE="${2:-cmd/server/main.go}"

if ! command -v swag >/dev/null 2>&1; then
  echo "swag command not found; please install github.com/swaggo/swag/cmd/swag first" >&2
  exit 1
fi

cd "${SERVICE_ROOT}"
mkdir -p api

swag init \
  -g "${ENTRY_FILE}" \
  -o "./api" \
  --outputTypes json

echo "generated ${SERVICE_ROOT}/api/swagger.json"
