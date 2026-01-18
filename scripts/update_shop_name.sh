#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8081}"
TOKEN="${TOKEN:-4f125923-4462-4967-b4cd-65b57fb98938}"
SHOP_ID="${SHOP_ID:-1}"
SHOP_NAME="${SHOP_NAME:-103茶餐厅-test222}"

url="${BASE_URL}/shop"
payload=$(printf '{"id":%s,"name":"%s"}' "${SHOP_ID}" "${SHOP_NAME}")

curl -sS -X PUT "${url}" \
  -H "authorization: ${TOKEN}" \
  -H "content-type: application/json" \
  -d "${payload}" | sed 's/$/\n/' || true
