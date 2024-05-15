package cli

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/app"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var (
	src, srcURL, srcServiceAccountID, srcAPIKey, dst, dstServiceAccountID string
	verbose                                                               bool
)

func migrate(src, dst app.GrafanaInput, verbose bool) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// create clients
	srcAWSClient := aws.New(sess, src.Region, src.IsGamma)
	srcGrafanaClient, err := src.CreateGrafanaAPIClient(srcAWSClient)
	if err != nil {
		return err
	}
	defer src.DeleteGrafanaAuth(srcAWSClient, srcGrafanaClient.Auth)

	dstAWSClient := aws.New(sess, dst.Region, dst.IsGamma)
	dstGrafanaClient, err := dst.CreateGrafanaAPIClient(dstAWSClient)
	if err != nil {
		return err
	}
	defer dst.DeleteGrafanaAuth(dstAWSClient, dstGrafanaClient.Auth)

	migrate := app.App{Src: srcGrafanaClient.Client, Dst: dstGrafanaClient.Client, Verbose: verbose}
	return migrate.Run()
}

// BuildMigrateCmd builds the migrate CLI command
func BuildMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Discover Managed Grafana workspaces",
		Long:  "Discover Managed Grafana workspaces",
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			log.Info()
			src, err := app.NewGrafanaInput(src, srcURL, srcServiceAccountID, srcAPIKey)
			if err != nil {
				return err
			}
			dst, err := app.NewGrafanaInput(dst, "", dstServiceAccountID, "")
			if err != nil {
				return err
			}
			return migrate(src, dst, verbose)
		}),
	}

	cmd.Flags().StringVarP(&src, "src", "s", "", "Source Grafana workspace")
	cmd.Flags().StringVarP(&srcServiceAccountID, "src-service-account-id", "", "", "Grafana Service Account ID for source workspace (exclusive with src)")
	cmd.Flags().StringVarP(&srcURL, "src-url", "", "", "Source Grafana URL (exclusive with src)")
	cmd.Flags().StringVarP(&srcAPIKey, "src-api-key", "", "", "Source Grafana API Key or Service Account Token (mandatory when using src-url)")
	cmd.MarkFlagsRequiredTogether("src-url", "src-api-key")
	cmd.MarkFlagsMutuallyExclusive("src-url", "src")

	cmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination Grafana Workspace endpoint")
	cmd.Flags().StringVarP(&dstServiceAccountID, "dst-service-account-id", "", "", "Grafana Service Account ID for destination workspace (required for v10+ workspaces)")
	cmd.MarkFlagRequired("dst")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	return cmd
}
