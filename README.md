# Amazon Managed Grafana Migrator

Amazon Managed Grafana Migrator is a migration utility to help migrate Grafana's
content (data sources, dashboards, folders and alert rules) to Amazon Managed
Grafana.

- Migrating from and to Amazon Managed Grafana Workspace (copy workspace contents)
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

```bash
amazon-managed-grafana-migrator migrate \
  --src-url https://grafana.example.com/
  --src-api-key API_KEY_HERE
  --dst g-abcdef5678.grafana-workspace.us-west-2.amazonaws.com
```

### Getting help

```bash
amazon-managed-grafana-migrator --help

# command specific help
amazon-managed-grafana-migrator migrate --help
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for more information.

## License

This project is licensed under the Apache-2.0 License.

