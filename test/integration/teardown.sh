#!/usr/bin/env bash
# Teardown: clean migrated content from the destination AMG workspace
# Also tears down Docker Grafana if running.
#
# Required env vars:
#   DST_ENDPOINT              - AMG v12 workspace endpoint
#   DST_SERVICE_ACCOUNT_ID    - Service account ID for v12 workspace
#
# The script creates a temporary SA token, deletes all non-default content, then cleans up.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
DOCKER_DIR="$SCRIPT_DIR/docker"

info() { echo "→ $1"; }

# Stop Docker Grafana if running
if [ -f "$DOCKER_DIR/docker-compose.yml" ]; then
    info "Stopping Docker Grafana..."
    cd "$DOCKER_DIR" && docker compose down 2>/dev/null || true
fi

# Clean destination workspace
if [ -z "${DST_ENDPOINT:-}" ] || [ -z "${DST_SERVICE_ACCOUNT_ID:-}" ]; then
    echo "DST_ENDPOINT and DST_SERVICE_ACCOUNT_ID not set, skipping AMG cleanup"
    exit 0
fi

DST_URL="https://$DST_ENDPOINT"

# Create a temporary SA token via AWS CLI
info "Creating temporary SA token for cleanup..."
TOKEN_RESP=$(aws grafana create-workspace-service-account-token \
    --workspace-id "$(echo "$DST_ENDPOINT" | cut -d. -f1)" \
    --service-account-id "$DST_SERVICE_ACCOUNT_ID" \
    --name "teardown-$(date +%s)" \
    --seconds-to-live 300)

TOKEN=$(echo "$TOKEN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['serviceAccountToken']['key'])")
TOKEN_ID=$(echo "$TOKEN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['serviceAccountToken']['id'])")

cleanup_token() {
    info "Removing temporary SA token..."
    aws grafana delete-workspace-service-account-token \
        --workspace-id "$(echo "$DST_ENDPOINT" | cut -d. -f1)" \
        --service-account-id "$DST_SERVICE_ACCOUNT_ID" \
        --token-id "$TOKEN_ID" 2>/dev/null || true
}
trap cleanup_token EXIT

AUTH="Authorization: Bearer $TOKEN"

# Delete alert rules
info "Deleting alert rules..."
RULES=$(curl -sf -H "$AUTH" "$DST_URL/api/ruler/grafana/api/v1/rules" 2>/dev/null || echo "{}")
for folder in $(echo "$RULES" | grep -o '"[^"]*":' | tr -d '":' | head -20); do
    RULE_UIDS=$(curl -sf -H "$AUTH" "$DST_URL/api/ruler/grafana/api/v1/rules/$folder" 2>/dev/null | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
    for uid in $RULE_UIDS; do
        curl -sf -H "$AUTH" -X DELETE "$DST_URL/api/v1/provisioning/alert-rules/$uid" > /dev/null 2>&1 || true
    done
done

# Delete dashboards
info "Deleting dashboards..."
DASH_UIDS=$(curl -sf -H "$AUTH" "$DST_URL/api/search?type=dash-db&limit=100" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
for uid in $DASH_UIDS; do
    curl -sf -H "$AUTH" -X DELETE "$DST_URL/api/dashboards/uid/$uid" > /dev/null 2>&1 || true
done

# Delete folders
info "Deleting folders..."
FOLDER_UIDS=$(curl -sf -H "$AUTH" "$DST_URL/api/folders?limit=100" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
for uid in $FOLDER_UIDS; do
    curl -sf -H "$AUTH" -X DELETE "$DST_URL/api/folders/$uid" > /dev/null 2>&1 || true
done

# Delete datasources
info "Deleting datasources..."
DS_UIDS=$(curl -sf -H "$AUTH" "$DST_URL/api/datasources" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
for uid in $DS_UIDS; do
    curl -sf -H "$AUTH" -X DELETE "$DST_URL/api/datasources/uid/$uid" > /dev/null 2>&1 || true
done

info "Teardown complete."
