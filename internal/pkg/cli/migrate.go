package cli

import (
	"amazon-managed-grafana-migrator/internal/pkg/app"
	"amazon-managed-grafana-migrator/internal/pkg/aws"
	"amazon-managed-grafana-migrator/internal/pkg/grafana"
	"amazon-managed-grafana-migrator/internal/pkg/log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var (
	src, srcURL, srcAPIKey, dst string
)

func migrate(src, dst app.GrafanaInput) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// create clients
	srcAWSClient := aws.New(sess, src.Region, src.IsGamma)
	srcGrafanaClient, err := src.CreateGrafanaAPIClient(srcAWSClient)
	if err != nil {
		return err
	}
	defer src.DeleteAPIKeys(srcAWSClient, srcGrafanaClient.Key)

	dstAWSClient := aws.New(sess, dst.Region, dst.IsGamma)
	dstGrafanaClient, err := dst.CreateGrafanaAPIClient(dstAWSClient)
	if err != nil {
		return err
	}
	defer dst.DeleteAPIKeys(dstAWSClient, dstGrafanaClient.Key)

	//looking for API key for CLI provided or AWS CLI SDK
	apikey := srcGrafanaClient.Key.APIKey
	if apikey == "" {
		apikey = srcGrafanaClient.Input.APIKey
	}

	// new custom client
	customClient, err := grafana.New(srcGrafanaClient.Input.URL, apikey)
	if err != nil {
		return err
	}

	migrate := app.App{Src: srcGrafanaClient.Client, Dst: dstGrafanaClient.Client}
	return migrate.Run(app.CustomGrafanaClient{Client: customClient})
}

// BuildMigrateCmd builds the migrate CLI command
func BuildMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Discover Managed Grafana workspaces",
		Long:  "Discover Managed Grafana workspaces",
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			log.Info()
			src, err := app.NewGrafanaInput(src, srcURL, srcAPIKey)
			if err != nil {
				return err
			}
			dst, err := app.NewGrafanaInput(dst, "", "")
			if err != nil {
				return err
			}
			return migrate(src, dst)
		}),
	}

	cmd.Flags().StringVarP(&src, "src", "s", "", "Source Grafana workspace")
	cmd.Flags().StringVarP(&srcURL, "src-url", "", "", "Source Grafana URL (exclusive with src)")
	cmd.Flags().StringVarP(&srcAPIKey, "src-api-key", "", "", "Source Grafana API Key (mandatory when using src-url)")
	cmd.MarkFlagsRequiredTogether("src-url", "src-api-key")
	cmd.MarkFlagsMutuallyExclusive("src-url", "src")

	cmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination Grafana Workspace endpoint")
	cmd.MarkFlagRequired("dst")
	return cmd
}
