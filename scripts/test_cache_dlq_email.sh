#!/usr/bin/env bash
set -euo pipefail

REDIS_CONTAINER="${REDIS_CONTAINER:-redis-hmdp}"
LOG_FILE="${LOG_FILE:-server.log}"

cleanup() {
  docker start "${REDIS_CONTAINER}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "stopping redis: ${REDIS_CONTAINER}"
docker stop "${REDIS_CONTAINER}" >/dev/null

echo "triggering shop update"
./scripts/update_shop_name.sh

echo "restoring redis: ${REDIS_CONTAINER}"
docker start "${REDIS_CONTAINER}" >/dev/null

echo "checking dlq/email logs"
if command -v rg >/dev/null 2>&1; then
  rg -n "shop cache delete failed|cache invalidate delete failed|cache invalidate dlq|cache invalidate dlq email" "${LOG_FILE}" || true
else
  grep -n "shop cache delete failed\|cache invalidate delete failed\|cache invalidate dlq\|cache invalidate dlq email" "${LOG_FILE}" || true
fi
