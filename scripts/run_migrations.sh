#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
. "$BASE_DIR/.env"

PATH="$(go env GOPATH)/bin:$PATH"

export PATH

PG_URL=${PG_URL:-"postgres://user:myAwEsOm3pa55@w0rd@localhost:5432/db?sslmode=disable"}

echo "running migrations with PG_URL=$PG_URL"

migrate -path "$BASE_DIR/migrations" -database "$PG_URL" up
