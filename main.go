// Main CLI entrypoint.
package main

import (
	"os"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/cli"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"
	"github.com/spf13/cobra"
)

const (
	shortDescription = "Amazon Managed Grafana migration utility"
	version          = "0.1.3"
)

var region string

func main() {

	cmd := buildRootCmd()
	if err := cmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func buildRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "amazon-managed-grafana-migrator",
		Short: shortDescription,
		Example: `
Displays the help menu for the "migrate" command.
$ amazon-managed-grafana-migrator migrate --help

Discovers all the workspaces in the specified region.
$ amazon-managed-grafana-migrator discover --region eu-west-1
		`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Version = version
	cmd.AddCommand(cli.BuildDiscoverCmd())
	cmd.AddCommand(cli.BuildMigrateCmd())
	cmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS Region")

	return cmd
}
