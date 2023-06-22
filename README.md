# Amazon Managed Grafana Migrator

[![Build](https://github.com/aws-observability/amazon-managed-grafana-migrator/actions/workflows/go.yml/badge.svg)](https://github.com/aws-observability/amazon-managed-grafana-migrator/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aws-observability/amazon-managed-grafana-migrator)](https://goreportcard.com/report/github.com/aws-observability/amazon-managed-grafana-migrator)

Amazon Managed Grafana Migrator is a CLI migration utility to migrate Grafana
content (data sources, dashboards, folders and alert rules) to Amazon Managed
Grafana. It supports the following migration scenarios:

- Migrating from and to Amazon Managed Grafana Workspace (eg. Moving to v9.4)
- Migrating from a Grafana server to an Amazon Managed Grafana Workspace

<img src="https://user-images.githubusercontent.com/10175027/235176809-9b71af1a-79a9-416a-b26e-ccdf725779d7.gif" width="80%" height="80%"/>


:warning: Alerting rules migration are only supported in Amazon Managed Grafana
v9.4

## Installation

Build from latest release. This requires Go installed on your environement.

```console
go install github.com/aws-observability/amazon-managed-grafana-migrator@latest
```

This command above will build the binary into your Go path `($HOME/go/bin)`.
Make sure to add your Go bin in your $PATH to run the command.
For Linux, this is usually `export PATH=$PATH:$HOME/go/bin`.

You can also download the pre-compiled binary for your OS and CPU architecture
from our [GitHub releases](https://github.com/aws-observability/amazon-managed-grafana-migrator/releases/latest).

## Usage

### Discovering your Workspaces

```console
amazon-managed-grafana-migrator discover --region eu-west-1
```

### Migrating between Workspaces

```console
amazon-managed-grafana-migrator migrate \
  --src g-abcdef1234.grafana-workspace.eu-central-1.amazonaws.com \
  --dst g-abcdef5678.grafana-workspace.us-west-2.amazonaws.com
```

### Migrating to Amazon Managed Grafana

```console
amazon-managed-grafana-migrator migrate \
  --src-url https://grafana.example.com/
  --src-api-key API_KEY_HERE
  --dst g-abcdef5678.grafana-workspace.us-west-2.amazonaws.com
```

### Getting help

```console
amazon-managed-grafana-migrator --help

# command specific help
amazon-managed-grafana-migrator migrate --help
```

## Permissions

To run this tool you need AWS permissions through IAM. Make sure you have the
AWS command line tool installed and have already run `aws configure` before you
start. Below are the minimum permissions required by the tool:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "grafana:DeleteWorkspaceApiKey",
                "grafana:DescribeWorkspace",
                "grafana:CreateWorkspaceApiKey"
            ],
            "Resource": "arn:aws:grafana:*:<ACCOUNT_ID>:/workspaces/<WORKSPACE_ID>"
        },
        {
            "Effect": "Allow",
            "Action": "grafana:ListWorkspaces",
            "Resource": "*"
        }
    ]
}
```

If you a migrating from a Grafana server, you will need an active API Key with
"ADMIN" role.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for more information.


## License

This project is licensed under the Apache-2.0 License.
