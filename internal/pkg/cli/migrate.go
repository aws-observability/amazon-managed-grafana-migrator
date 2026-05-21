package cli

import (
	"context"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/app"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var (
	src, srcURL, srcServiceAccountID, srcAPIKey, dst, dstServiceAccountID string
	verbose                                                               bool
)

func migrate(srcInput, dstInput app.GrafanaInput, verbose bool) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// create source client
	srcAWSClient := aws.New(cfg, srcInput.Region)
	srcGrafana, err := srcInput.CreateGrafanaHTTPClient(ctx, srcAWSClient)
	if err != nil {
		return err
	}
	defer srcInput.DeleteGrafanaAuth(ctx, srcAWSClient, srcGrafana.Auth)

	// create destination client
	dstAWSClient := aws.New(cfg, dstInput.Region)
	dstGrafana, err := dstInput.CreateGrafanaHTTPClient(ctx, dstAWSClient)
	if err != nil {
		return err
	}
	defer dstInput.DeleteGrafanaAuth(ctx, dstAWSClient, dstGrafana.Auth)

	migrator := app.App{
		Src:     app.NewGrafanaClient(srcGrafana.BaseURL, srcGrafana.Auth.GetAuth(), srcGrafana.HTTPClient),
		Dst:     app.NewGrafanaClient(dstGrafana.BaseURL, dstGrafana.Auth.GetAuth(), dstGrafana.HTTPClient),
		Verbose: verbose,
	}
	return migrator.Run()
}

// BuildMigrateCmd builds the migrate CLI command
func BuildMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate Grafana content between workspaces",
		Long:  "Migrate data sources, dashboards, folders, and alert rules between Grafana workspaces",
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			log.Info()
			srcInput, err := app.NewGrafanaInput(src, srcURL, srcServiceAccountID, srcAPIKey)
			if err != nil {
				return err
			}
			dstInput, err := app.NewGrafanaInput(dst, "", dstServiceAccountID, "")
			if err != nil {
				return err
			}
			return migrate(srcInput, dstInput, verbose)
		}),
	}

	cmd.Flags().StringVarP(&src, "src", "s", "", "Source AMG workspace endpoint (e.g. g-xxx.grafana-workspace.us-east-1.amazonaws.com)")
	cmd.Flags().StringVarP(&srcServiceAccountID, "src-service-account-id", "", "", "Grafana Service Account ID for source workspace")
	cmd.Flags().StringVarP(&srcURL, "src-url", "", "", "Source Grafana URL (for non-AMG Grafana servers)")
	cmd.Flags().StringVarP(&srcAPIKey, "src-api-key", "", "", "Source Grafana API Key (required with --src-url)")
	cmd.MarkFlagsRequiredTogether("src-url", "src-api-key")
	cmd.MarkFlagsMutuallyExclusive("src-url", "src")

	cmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination AMG workspace endpoint")
	cmd.Flags().StringVarP(&dstServiceAccountID, "dst-service-account-id", "", "", "Grafana Service Account ID for destination workspace (required)")
	cmd.MarkFlagRequired("dst")
	cmd.MarkFlagRequired("dst-service-account-id")

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	return cmd
}
