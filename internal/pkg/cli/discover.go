// Package cli provides the CLI
package cli

import (
	"errors"
	"os"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var errMissingRegion = errors.New("missing AWS region")

// discover discovers workspaces in a region
func discover(region string) error {
	if region == "" {
		return errMissingRegion
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	awsGrafana := aws.New(sess, region, false)

	wx, err := awsGrafana.ListWorkspaces()

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

// runCmdE wraps one of the run error methods, PreRunE, RunE, of a cobra command so that if a user
// types "help" in the arguments the usage string is printed instead of running the command.
func runCmdE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "help" {
			_ = cmd.Help() // Help always returns nil.
			os.Exit(0)
		}
		return f(cmd, args)
	}
}
