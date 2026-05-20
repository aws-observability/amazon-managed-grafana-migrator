// Package cli provides the CLI
package cli

import (
	"context"
	"errors"
	"os"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var errMissingRegion = errors.New("missing AWS region")

func discover(region string) error {
	if region == "" {
		return errMissingRegion
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	awsGrafana := aws.New(cfg, region)
	wx, err := awsGrafana.ListWorkspaces(ctx)
	if err != nil {
		return err
	}

	if len(wx) == 0 {
		log.Info("No workspaces found")
	} else {
		log.Success("Discovered ", len(wx), " workspaces")
		log.Info()
	}

	for _, w := range wx {
		log.Infof("Version: %s\nName: %s\nEndpoint: %s\n\n", w.Version, w.Name, w.Endpoint)
	}
	return nil
}

// BuildDiscoverCmd builds the command for discovering workspaces
func BuildDiscoverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover Managed Grafana workspaces",
		Long:  "Discover Managed Grafana workspaces in a region",
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			if region == "" {
				return errMissingRegion
			}
			return discover(region)
		}),
	}
	return cmd
}

// runCmdE wraps a cobra RunE so that "help" as an argument prints usage
func runCmdE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "help" {
			_ = cmd.Help()
			os.Exit(0)
		}
		return f(cmd, args)
	}
}
