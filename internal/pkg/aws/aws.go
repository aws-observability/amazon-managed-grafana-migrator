// Package aws provides a wrapper around the AWS Grafana API
package aws

import (
	"fmt"
	"time"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
)

type api interface {
	ListWorkspaces(input *managedgrafana.ListWorkspacesInput) (*managedgrafana.ListWorkspacesOutput, error)
	CreateWorkspaceApiKey(input *managedgrafana.CreateWorkspaceApiKeyInput) (*managedgrafana.CreateWorkspaceApiKeyOutput, error)
	DeleteWorkspaceApiKey(input *managedgrafana.DeleteWorkspaceApiKeyInput) (*managedgrafana.DeleteWorkspaceApiKeyOutput, error)
	//DescribeWorkspace(input *managedgrafana.DescribeWorkspaceInput) (*managedgrafana.DescribeWorkspaceOutput, error)
}

// AMGApiKey holds the dataplane Grafana API key for a workspace
type AMGApiKey struct {
	KeyName, APIKey, WorkspaceID string
}

// AMG is a AWS SDK client for AMG apis
type AMG struct {
	Client api
}

// Workspace contains informations about a Grafana workspace
type Workspace struct {
	ID       string
	Name     string
	Version  string
	Endpoint string
	Region   string
	APIKey   AMGApiKey
}

// New creates a new AMG client
func New(s *session.Session, region string, isGamma bool) *AMG {
	if isGamma {
		return &AMG{
			Client: managedgrafana.New(s, aws.NewConfig().WithRegion(region).WithEndpoint("https://grafana-gamma.us-east-1.amazonaws.com")),
		}
	}
	return &AMG{
		Client: managedgrafana.New(s, aws.NewConfig().WithRegion(region)),
	}
}

// ListWorkspaces lists all workspaces in the region
func (a *AMG) ListWorkspaces() ([]Workspace, error) {
	response, err := a.Client.ListWorkspaces(&managedgrafana.ListWorkspacesInput{})

	if err != nil {
		return []Workspace{}, err
	}

	wx := []Workspace{}

	for _, workspace := range response.Workspaces {
		w := Workspace{
			ID:       *workspace.Id,
			Name:     *workspace.Name,
			Version:  *workspace.GrafanaVersion,
			Endpoint: *workspace.Endpoint,
		}
		wx = append(wx, w)
	}
	return wx, nil
}

// CreateWorkspaceApiKey creates a new API key for a workspace
//
//revive:disable
func (a *AMG) CreateWorkspaceApiKey(workspaceID string) (AMGApiKey, error) {
	currentTime := time.Now().UTC()
	keyName := fmt.Sprintf("%s-%d", "amg-migrator", currentTime.UnixMilli())

	log.Debug("Creating temporary API key for ", workspaceID)

	duration := time.Duration(5 * 60 * time.Second)
	resp, err := a.Client.CreateWorkspaceApiKey(&managedgrafana.CreateWorkspaceApiKeyInput{
		KeyName:       aws.String(keyName),
		KeyRole:       aws.String("ADMIN"),
		SecondsToLive: aws.Int64(int64(duration.Seconds())),
		WorkspaceId:   aws.String(workspaceID),
	})

	if err != nil {
		return AMGApiKey{}, err
	}

	return AMGApiKey{keyName, *resp.Key, workspaceID}, nil
}

// DeleteWorkspaceApiKey deletes an API key for a workspace
func (a *AMG) DeleteWorkspaceApiKey(apiKey AMGApiKey) error {
	log.Info()
	log.Debug("Removing temporary API key for ", apiKey.WorkspaceID)

	_, err := a.Client.DeleteWorkspaceApiKey(&managedgrafana.DeleteWorkspaceApiKeyInput{
		KeyName:     aws.String(apiKey.KeyName),
		WorkspaceId: aws.String(apiKey.WorkspaceID),
	})

	if err != nil {
		log.Error(err)
	}
	return err
}
