#!/usr/bin/env bash
# Integration test runner for amazon-managed-grafana-migrator
# Runs two scenarios:
#   1. Docker Grafana v10 → AMG v12 (requires Docker + AMG v12 workspace)
#   2. AMG v9 → AMG v12 (requires both AMG workspaces provisioned via Terraform)
#
# Usage:
#   ./run_integration.sh                    # Run all tests
#   ./run_integration.sh docker-to-amg      # Run only Docker→AMG scenario
#   ./run_integration.sh amg-to-amg         # Run only AMG v9→v12 scenario
#   ./run_integration.sh v10-to-v12         # Run only AMG v10→v12 scenario
#
# Required environment variables:
#   DST_ENDPOINT              - AMG v12 workspace endpoint
#   DST_SERVICE_ACCOUNT_ID    - Service account ID for v12 workspace
#
# For AMG v9→v12 scenario additionally:
#   SRC_ENDPOINT              - AMG v9 workspace endpoint
#   SRC_SERVICE_ACCOUNT_ID    - Service account ID for v9 workspace
#
# For AMG v10→v12 scenario additionally:
#   SRC_V10_ENDPOINT          - AMG v10 workspace endpoint
#   SRC_V10_SERVICE_ACCOUNT_ID - Service account ID for v10 workspace
#
# For Docker→AMG scenario:
#   Docker must be running (docker-compose will be started automatically)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
DOCKER_DIR="$SCRIPT_DIR/docker"
BINARY="$ROOT_DIR/amazon-managed-grafana-migrator"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

pass() { echo -e "${GREEN}✔ PASS:${NC} $1"; }
fail() { echo -e "${RED}✘ FAIL:${NC} $1"; exit 1; }
info() { echo -e "${YELLOW}→${NC} $1"; }

# Build the migrator binary
build_binary() {
    info "Building migrator binary..."
    cd "$ROOT_DIR"
    GOPROXY=direct go build -o "$BINARY" .
}

# Wait for Grafana to be healthy
wait_for_grafana() {
    local url="$1"
    local max_attempts=30
    info "Waiting for Grafana at $url..."
    for i in $(seq 1 $max_attempts); do
        if curl -sf "$url/api/health" > /dev/null 2>&1; then
            return 0
        fi
        sleep 2
    done
    fail "Grafana at $url did not become healthy"
}

# Create a service account + token on local Grafana (for Docker scenario)
create_local_sa_token() {
    local url="$1"
    local auth="admin:admin"

    # Create service account
    local sa_resp
    sa_resp=$(curl -sf -X POST "$url/api/serviceaccounts" \
        -u "$auth" \
        -H "Content-Type: application/json" \
        -d '{"name":"migrator-integ-test","role":"Admin"}' 2>/dev/null || true)

    local sa_id
    sa_id=$(echo "$sa_resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))" 2>/dev/null || true)

    if [ -z "$sa_id" ]; then
        # SA might already exist, try to find it
        sa_id=$(curl -sf "$url/api/serviceaccounts/search?query=migrator-integ-test" \
            -u "$auth" | python3 -c "import sys,json; r=json.load(sys.stdin); print(r['serviceAccounts'][0]['id'] if r.get('serviceAccounts') else '')" 2>/dev/null || true)
    fi

    # Create token
    local token_resp
    token_resp=$(curl -sf -X POST "$url/api/serviceaccounts/$sa_id/tokens" \
        -u "$auth" \
        -H "Content-Type: application/json" \
        -d '{"name":"integ-test-token"}')

    echo "$token_resp" | python3 -c "import sys,json; print(json.load(sys.stdin)['key'])"
}

# Verify migration results on destination
verify_migration() {
    local dst_url="$1"
    local token="$2"
    local scenario="$3"

    info "Verifying migration results for: $scenario"

    # Check datasources
    local ds_count
    ds_count=$(curl -sf -H "Authorization: Bearer $token" "$dst_url/api/datasources" | grep -o '"name"' | wc -l)
    if [ "$ds_count" -ge 1 ]; then
        pass "$scenario: Found $ds_count datasource(s)"
    else
        fail "$scenario: No datasources found in destination"
    fi

    # Check folders
    local folder_count
    folder_count=$(curl -sf -H "Authorization: Bearer $token" "$dst_url/api/folders" | grep -o '"title"' | wc -l)
    if [ "$folder_count" -ge 1 ]; then
        pass "$scenario: Found $folder_count folder(s)"
    else
        fail "$scenario: No folders found in destination"
    fi

    # Check dashboards
    local dash_count
    dash_count=$(curl -sf -H "Authorization: Bearer $token" "$dst_url/api/search?type=dash-db" | grep -o '"uid"' | wc -l)
    if [ "$dash_count" -ge 1 ]; then
        pass "$scenario: Found $dash_count dashboard(s)"
    else
        fail "$scenario: No dashboards found in destination"
    fi
}

# Scenario 1: Docker Grafana v10 → AMG v12
test_docker_to_amg() {
    info "=== Scenario: Docker Grafana v10 → AMG v12 ==="

    if [ -z "${DST_ENDPOINT:-}" ] || [ -z "${DST_SERVICE_ACCOUNT_ID:-}" ]; then
        fail "DST_ENDPOINT and DST_SERVICE_ACCOUNT_ID must be set"
    fi

    # Start local Grafana
    info "Starting local Grafana v10..."
    cd "$DOCKER_DIR"
    docker compose up -d
    wait_for_grafana "http://localhost:3000"

    # Get API token from local Grafana
    local src_token
    src_token=$(create_local_sa_token "http://localhost:3000")
    if [ -z "$src_token" ]; then
        fail "Could not create service account token on local Grafana"
    fi
    pass "Created local Grafana service account token"

    # Run migration
    info "Running migration: Docker Grafana v10 → AMG v12..."
    "$BINARY" migrate \
        --src-url "http://localhost:3000" \
        --src-api-key "$src_token" \
        --dst "$DST_ENDPOINT" \
        --dst-service-account-id "$DST_SERVICE_ACCOUNT_ID" \
        --verbose

    pass "Docker→AMG migration completed"

    # Cleanup local Grafana
    docker compose down
}

# Scenario 2: AMG v9 → AMG v12
test_amg_to_amg() {
    info "=== Scenario: AMG v9 → AMG v12 ==="

    if [ -z "${SRC_ENDPOINT:-}" ] || [ -z "${SRC_SERVICE_ACCOUNT_ID:-}" ]; then
        fail "SRC_ENDPOINT and SRC_SERVICE_ACCOUNT_ID must be set"
    fi
    if [ -z "${DST_ENDPOINT:-}" ] || [ -z "${DST_SERVICE_ACCOUNT_ID:-}" ]; then
        fail "DST_ENDPOINT and DST_SERVICE_ACCOUNT_ID must be set"
    fi

    info "Running migration: AMG v9 → AMG v12..."
    "$BINARY" migrate \
        --src "$SRC_ENDPOINT" \
        --src-service-account-id "$SRC_SERVICE_ACCOUNT_ID" \
        --dst "$DST_ENDPOINT" \
        --dst-service-account-id "$DST_SERVICE_ACCOUNT_ID" \
        --verbose

    pass "AMG v9→v12 migration completed"
}

# Scenario 3: AMG v10 → AMG v12
test_v10_to_v12() {
    info "=== Scenario: AMG v10 → AMG v12 ==="

    if [ -z "${SRC_V10_ENDPOINT:-}" ] || [ -z "${SRC_V10_SERVICE_ACCOUNT_ID:-}" ]; then
        fail "SRC_V10_ENDPOINT and SRC_V10_SERVICE_ACCOUNT_ID must be set"
    fi
    if [ -z "${DST_ENDPOINT:-}" ] || [ -z "${DST_SERVICE_ACCOUNT_ID:-}" ]; then
        fail "DST_ENDPOINT and DST_SERVICE_ACCOUNT_ID must be set"
    fi

    info "Running migration: AMG v10 → AMG v12..."
    "$BINARY" migrate \
        --src "$SRC_V10_ENDPOINT" \
        --src-service-account-id "$SRC_V10_SERVICE_ACCOUNT_ID" \
        --dst "$DST_ENDPOINT" \
        --dst-service-account-id "$DST_SERVICE_ACCOUNT_ID" \
        --verbose

    pass "AMG v10→v12 migration completed"
}

# Main
main() {
    local scenario="${1:-all}"

    build_binary

    case "$scenario" in
        docker-to-amg)
            test_docker_to_amg
            ;;
        amg-to-amg)
            test_amg_to_amg
            ;;
        v10-to-v12)
            test_v10_to_v12
            ;;
        all)
            test_docker_to_amg
            test_amg_to_amg
            test_v10_to_v12
            ;;
        *)
            echo "Usage: $0 [docker-to-amg|amg-to-amg|v10-to-v12|all]"
            exit 1
            ;;
    esac

    echo ""
    pass "All integration tests passed!"
}

main "$@"
