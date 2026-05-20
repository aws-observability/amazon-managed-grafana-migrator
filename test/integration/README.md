# Integration Tests

## Overview

Two migration scenarios are tested:

| Scenario | Source | Destination |
|----------|--------|-------------|
| `docker-to-amg` | Local Grafana v10.4 (Docker) | AMG v12.4 |
| `amg-to-amg` | AMG v9.4 | AMG v12.4 |

## Prerequisites

- AWS credentials configured (`aws configure` or env vars)
- Docker (for `docker-to-amg` scenario)
- Go 1.22+

## Setup: Provision AMG Workspaces

```bash
cd terraform/
terraform init
terraform apply
```

This creates:
- AMG v9.4 workspace (source) with an ADMIN service account
- AMG v12.4 workspace (destination) with an ADMIN service account

Export the outputs:

```bash
export SRC_ENDPOINT=$(terraform output -raw src_v9_endpoint)
export SRC_SERVICE_ACCOUNT_ID=$(terraform output -raw src_v9_service_account_id)
export DST_ENDPOINT=$(terraform output -raw dst_v12_endpoint)
export DST_SERVICE_ACCOUNT_ID=$(terraform output -raw dst_v12_service_account_id)
```

## Seed Source Workspace

After provisioning, seed the v9 workspace with test content:

```bash
./seed.sh
```

This creates 2 datasources, 1 folder, 1 dashboard (2 panels), and 1 alert rule
in the source v9 workspace using a temporary service account token.

## Running Tests

```bash
# Run all scenarios
./run_integration.sh

# Run only Docker→AMG
./run_integration.sh docker-to-amg

# Run only AMG→AMG (requires seed data in v9 workspace)
./run_integration.sh amg-to-amg
```

## Cleanup

```bash
# Remove migrated content from destination + stop Docker
./teardown.sh

# Destroy AWS infrastructure
cd terraform/ && terraform destroy
```

## CI

The integration tests run via GitHub Actions with `workflow_dispatch` (manual trigger).
Required secrets:
- `INTEG_AWS_ROLE_ARN` — IAM role ARN for OIDC auth
- `INTEG_DST_ENDPOINT` / `INTEG_DST_SERVICE_ACCOUNT_ID`
- `INTEG_SRC_ENDPOINT` / `INTEG_SRC_SERVICE_ACCOUNT_ID`
