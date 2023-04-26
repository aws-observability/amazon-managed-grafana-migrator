# Amazon Managed Grafana Migrator

[![Build](https://github.com/aws-observability/amazon-managed-grafana-migrator/actions/workflows/go.yml/badge.svg)](https://github.com/aws-observability/amazon-managed-grafana-migrator/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aws-observability/amazon-managed-grafana-migrator)](https://goreportcard.com/report/github.com/aws-observability/amazon-managed-grafana-migrator)

Amazon Managed Grafana Migrator is a CLI migration utility to migrate Grafana
content (data sources, dashboards, folders and alert rules) to Amazon Managed
Grafana. It supports the following migration scenarios:

- Migrating from and to Amazon Managed Grafana Workspace (eg. Moving to v9.4)
- Migrating from a Grafana server to an Amazon Managed Grafana Workspace

:warning: Alerting rules migration are only supported in Amazon Managed Grafana
v9.4


## Build

```console
go install github.com/aws-observability/amazon-managed-grafana-migrator@latest
```


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


## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for more information.


## License

This project is licensed under the Apache-2.0 License.
