# Amazon Managed Grafana Migrator

[![Build](https://github.com/aws-observability/amazon-managed-grafana-migrator/actions/workflows/go.yml/badge.svg)](https://github.com/aws-observability/amazon-managed-grafana-migrator/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aws-observability/amazon-managed-grafana-migrator)](https://goreportcard.com/report/github.com/aws-observability/amazon-managed-grafana-migrator)

🎉 May-15-24: v0.2.0 supports Grafana Service Accounts for v9 and v10 workspaces.
See [Amazon Managed Grafana announces support for Grafana version 10.4]()

🚨 Jul-25-23: Alerts migration are currently disabled with v0.1.9. See [#19](https://github.com/aws-observability/amazon-managed-grafana-migrator/issues/19)

🎉 Jul-19-23: [Amazon Grafana supports now in-place update from v8.4 to v9.4](https://aws.amazon.com/about-aws/whats-new/2023/07/amazon-managed-grafana-in-place-update-version-9-4/)

Amazon Managed Grafana Migrator is a CLI migration utility to migrate Grafana
content (data sources, dashboards, folders and alert rules) to Amazon Managed
Grafana. It supports the following migration scenarios:

- Migrating from and to Amazon Managed Grafana Workspace (eg. Moving to v10.4), although consider using the native functionality in the AWS Console, after testing
- Migrating from a Grafana server to an Amazon Managed Grafana Workspace

<img src="https://user-images.githubusercontent.com/10175027/235176809-9b71af1a-79a9-416a-b26e-ccdf725779d7.gif" width="80%" height="80%"/>

Amazon Managed Grafana v10.4 workspaces will require to
provide an `ADMIN` level [Grafana Service Account]() with the
`--src-service-account-id` or `--src-service-account-id` flags.

## Installation

Build from latest release. This requires Go (1.21 +) installed on your environement.

```console
go install github.com/aws-observability/amazon-managed-grafana-migrator@latest
```

This command above will build the binary into your Go path `($HOME/go/bin)`.
Make sure to add your Go bin in your $PATH to run the command.
For Linux, this is usually `export PATH=$PATH:$HOME/go/bin`.

You can also download the pre-compiled binary for your OS and CPU architecture
from our [GitHub releases](https://github.com/aws-observability/amazon-managed-grafana-migrator/releases/latest).

Example on Amazon Linux:

```console
wget https://github.com/aws-observability/amazon-managed-grafana-migrator/releases/download/v0.1.11/amazon-managed-grafana-migrator-linux-amd64.tar.gz
tar -zxvf amazon-managed-grafana-migrator-linux-amd64.tar.gz
sudo mv amazon-managed-grafana-migrator /usr/local/bin/
amazon-managed-grafana-migrator -v
```

## Usage

### Discovering your Workspaces

```console
amazon-managed-grafana-migrator discover --region eu-west-1
```

### Migrating to Amazon Managed Grafana v10

v9 and v10 introduced Grafana Service Accounts which will be required by the
migrator, especially for v10. Note that Service Accounts are billed as active
users

1. Creating a Service Account

```console
 aws grafana create-workspace-service-account --workspace-id g-abcdef5678 \
    --grafana-role ADMIN \
    --name <SA_NAME_HERE>
```

2. Running the migration

```console
amazon-managed-grafana-migrator migrate \
  --src-url https://grafana.example.com/
  --src-api-key API_KEY_HERE
  --dst g-abcdef5678.grafana-workspace.us-west-2.amazonaws.com
  --dst-service-account-id SERVICE_ACCOUNT_ID_HERE
```

### Migrating between Workspaces

```console
amazon-managed-grafana-migrator migrate \
  --src g-abcdef1234.grafana-workspace.eu-central-1.amazonaws.com \
  --dst g-abcdef5678.grafana-workspace.us-west-2.amazonaws.com
```

Or for v9+ workspaces:

```console
amazon-managed-grafana-migrator migrate \
  --src g-abcdef1234.grafana-workspace.eu-central-1.amazonaws.com \
  --src-service-account-id SERVICE_ACCOUNT_ID_HERE
  --dst g-abcdef5678.grafana-workspace.us-west-2.amazonaws.com
  --dst-service-account-id SERVICE_ACCOUNT_ID_HERE
```

### Migrating to Amazon Managed Grafana v8/v9

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
                "grafana:DescribeWorkspace",
                "grafana:CreateWorkspaceApiKey",
                "grafana:DeleteWorkspaceApiKey",
                "grafana:CreateWorkspaceServiceAccountToken",
                "grafana:DeleteWorkspaceServiceAccountToken"
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
