#!/usr/bin/env bash
set -euo pipefail

echo "bootstrap aevons gateway"
echo "1. start APISIX: cd apisix && docker compose up -d"
echo "2. start console: cd ../console && go run ./cmd/server"
