#!/usr/bin/env bash
# Seeds the source AMG v9 workspace with test content (datasources, folders, dashboards, alert rules).
# Run AFTER terraform apply. Uses a temporary service account token.
#
# Required env vars:
#   SRC_ENDPOINT              - AMG v9 workspace endpoint
#   SRC_SERVICE_ACCOUNT_ID    - Service account ID for v9 workspace

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

info() { echo "→ $1"; }
pass() { echo "✔ $1"; }

if [ -z "${SRC_ENDPOINT:-}" ] || [ -z "${SRC_SERVICE_ACCOUNT_ID:-}" ]; then
    echo "Error: SRC_ENDPOINT and SRC_SERVICE_ACCOUNT_ID must be set"
    exit 1
fi

WORKSPACE_ID=$(echo "$SRC_ENDPOINT" | cut -d. -f1)
SRC_URL="https://$SRC_ENDPOINT"

# Create temporary SA token
info "Creating temporary SA token..."
TOKEN_RESP=$(aws grafana create-workspace-service-account-token \
    --workspace-id "$WORKSPACE_ID" \
    --service-account-id "$SRC_SERVICE_ACCOUNT_ID" \
    --name "seed-$(date +%s)" \
    --seconds-to-live 600)

TOKEN=$(echo "$TOKEN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['serviceAccountToken']['key'])")
TOKEN_ID=$(echo "$TOKEN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['serviceAccountToken']['id'])")

cleanup_token() {
    info "Removing temporary SA token..."
    aws grafana delete-workspace-service-account-token \
        --workspace-id "$WORKSPACE_ID" \
        --service-account-id "$SRC_SERVICE_ACCOUNT_ID" \
        --token-id "$TOKEN_ID" 2>/dev/null || true
}
trap cleanup_token EXIT

AUTH="Authorization: Bearer $TOKEN"

# Create datasources
info "Creating datasources..."
curl -sf -X POST "$SRC_URL/api/datasources" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{
        "name": "TestPrometheus",
        "type": "prometheus",
        "access": "proxy",
        "url": "http://localhost:9090",
        "isDefault": true
    }' > /dev/null
pass "Created TestPrometheus datasource"

curl -sf -X POST "$SRC_URL/api/datasources" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{
        "name": "TestCloudWatch",
        "type": "cloudwatch",
        "access": "proxy",
        "jsonData": {"authType": "default", "defaultRegion": "us-east-1"}
    }' > /dev/null
pass "Created TestCloudWatch datasource"

# Create folder
info "Creating folder..."
FOLDER_RESP=$(curl -sf -X POST "$SRC_URL/api/folders" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"title": "Integration Test Folder", "uid": "integ-test-folder"}')
FOLDER_UID=$(echo "$FOLDER_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['uid'])")
pass "Created folder: $FOLDER_UID"

# Create dashboard
info "Creating dashboard..."
curl -sf -X POST "$SRC_URL/api/dashboards/db" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{
        "dashboard": {
            "uid": "integ-test-dash-1",
            "title": "Integration Test Dashboard",
            "tags": ["integration-test"],
            "panels": [
                {
                    "id": 1,
                    "title": "Sample Panel",
                    "type": "timeseries",
                    "datasource": {"type": "prometheus", "uid": "TestPrometheus"},
                    "targets": [{"expr": "up", "refId": "A"}],
                    "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0}
                },
                {
                    "id": 2,
                    "title": "Stat Panel",
                    "type": "stat",
                    "datasource": {"type": "prometheus", "uid": "TestPrometheus"},
                    "targets": [{"expr": "process_resident_memory_bytes", "refId": "A"}],
                    "gridPos": {"h": 4, "w": 6, "x": 12, "y": 0}
                }
            ],
            "schemaVersion": 39
        },
        "folderUid": "integ-test-folder",
        "overwrite": true
    }' > /dev/null
pass "Created dashboard: integ-test-dash-1"

# Create alert rule
info "Creating alert rule..."
# Get the prometheus datasource UID (may differ from name)
PROM_UID=$(curl -sf -H "$AUTH" "$SRC_URL/api/datasources" | python3 -c "
import sys, json
ds = json.load(sys.stdin)
for d in ds:
    if d['name'] == 'TestPrometheus':
        print(d['uid'])
        break
")

# Use ruler API (works on v9+) with folder title in URL path
curl -sf -X POST "$SRC_URL/api/ruler/grafana/api/v1/rules/Integration%20Test%20Folder" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d "{
        \"name\": \"integration-test-alerts\",
        \"interval\": \"1m\",
        \"rules\": [{
            \"grafana_alert\": {
                \"title\": \"Test Alert - Instance Down\",
                \"condition\": \"C\",
                \"no_data_state\": \"NoData\",
                \"exec_err_state\": \"Error\",
                \"data\": [
                    {
                        \"refId\": \"A\",
                        \"relativeTimeRange\": {\"from\": 600, \"to\": 0},
                        \"datasourceUid\": \"$PROM_UID\",
                        \"model\": {\"expr\": \"up == 0\", \"refId\": \"A\"}
                    },
                    {
                        \"refId\": \"C\",
                        \"relativeTimeRange\": {\"from\": 600, \"to\": 0},
                        \"datasourceUid\": \"__expr__\",
                        \"model\": {
                            \"type\": \"classic_conditions\",
                            \"refId\": \"C\",
                            \"conditions\": [{
                                \"evaluator\": {\"type\": \"gt\", \"params\": [0]},
                                \"operator\": {\"type\": \"and\"},
                                \"query\": {\"params\": [\"A\"]},
                                \"reducer\": {\"type\": \"last\"}
                            }]
                        }
                    }
                ]
            },
            \"for\": \"5m\",
            \"annotations\": {\"summary\": \"Test alert for integration testing\"},
            \"labels\": {\"severity\": \"warning\", \"team\": \"integration-test\"}
        }]
    }" > /dev/null
pass "Created alert rule"

echo ""
pass "Seed complete! Source workspace is ready for migration testing."
