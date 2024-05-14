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
	CreateWorkspaceServiceAccountToken(input *managedgrafana.CreateWorkspaceServiceAccountTokenInput) (*managedgrafana.CreateWorkspaceServiceAccountTokenOutput, error)
	DeleteWorkspaceApiKey(input *managedgrafana.DeleteWorkspaceApiKeyInput) (*managedgrafana.DeleteWorkspaceApiKeyOutput, error)
	DeleteWorkspaceServiceAccountToken(input *managedgrafana.DeleteWorkspaceServiceAccountTokenInput) (*managedgrafana.DeleteWorkspaceServiceAccountTokenOutput, error)
	DescribeWorkspace(input *managedgrafana.DescribeWorkspaceInput) (*managedgrafana.DescribeWorkspaceOutput, error)
}

type GrafanaAuth interface {
	GetAuth() string
}

// AMGApiKey is a GrafanaAuth struct for Grafana API key
type AMGApiKey struct {
	KeyName, APIKey, WorkspaceID string
}

// GetAuth returns the API key for a workspace
func (k AMGApiKey) GetAuth() string {
	return k.APIKey
}

// AMGServiceAccountToken is a GrafanaAuth struct for service account tokens
type AMGServiceAccountToken struct {
	ServiceAccountID, SATokenID, Token, WorkspaceID string
}

// GetAuth returns the SA token for a workspace
func (t AMGServiceAccountToken) GetAuth() string {
	return t.Token
}

// AMG is a AWS SDK client for AMG apis
type AMG struct {
	Client api
}

// Workspace contains information about a Grafana workspace
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
			ID:       aws.StringValue(workspace.Id),
			Name:     aws.StringValue(workspace.Name),
			Version:  aws.StringValue(workspace.GrafanaVersion),
			Endpoint: aws.StringValue(workspace.Endpoint),
		}
		wx = append(wx, w)
	}
	return wx, nil
}

// DescribeWorkspace returns information about a workspace
func (a *AMG) DescribeWorkspace(workspaceID string) (Workspace, error) {
	res, err := a.Client.DescribeWorkspace(&managedgrafana.DescribeWorkspaceInput{
		WorkspaceId: aws.String(workspaceID),
	})
	if err != nil {
		return Workspace{}, err
	}
	w := Workspace{
		ID:       workspaceID,
		Name:     aws.StringValue(res.Workspace.Name),
		Version:  aws.StringValue(res.Workspace.GrafanaVersion),
		Endpoint: aws.StringValue(res.Workspace.Endpoint),
	}
	return w, nil
}

// CreateWorkspaceApiKey creates a new API key for a workspace
//
//revive:disable
func (a *AMG) CreateWorkspaceApiKey(workspaceID string) (AMGApiKey, error) {
	currentTime := time.Now().UTC()
	keyName := fmt.Sprintf("%s-%d", "amg-migrator", currentTime.UnixMilli())

	log.InfoLight("Creating temporary API key for ", workspaceID)

	duration := time.Duration(30 * 60 * time.Second)
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
	log.InfoLight("Removing temporary API key for ", apiKey.WorkspaceID)

	_, err := a.Client.DeleteWorkspaceApiKey(&managedgrafana.DeleteWorkspaceApiKeyInput{
		KeyName:     aws.String(apiKey.KeyName),
		WorkspaceId: aws.String(apiKey.WorkspaceID),
	})

	if err != nil {
		log.Error(err)
	}
	return err
}

func (a *AMG) CreateServiceAccountToken(workspaceID, serviceAccountID string) (AMGServiceAccountToken, error) {
	log.Info()
	log.InfoLight("Creating service account token for service account ", serviceAccountID)
	currentTime := time.Now().UTC()
	saTokenName := fmt.Sprintf("%s-%d", "amg-migrator", currentTime.UnixMilli())

	duration := time.Duration(30 * 60 * time.Second)
	resp, err := a.Client.CreateWorkspaceServiceAccountToken(&managedgrafana.CreateWorkspaceServiceAccountTokenInput{
		Name:             aws.String(saTokenName),
		SecondsToLive:    aws.Int64(int64(duration.Seconds())),
		WorkspaceId:      aws.String(workspaceID),
		ServiceAccountId: aws.String(serviceAccountID),
	})

	if err != nil {
		return AMGServiceAccountToken{}, err
	}

	return AMGServiceAccountToken{
		ServiceAccountID: serviceAccountID,
		SATokenID:        *resp.ServiceAccountToken.Id,
		Token:            *resp.ServiceAccountToken.Key,
		WorkspaceID:      workspaceID,
	}, nil
}

func (a *AMG) DeleteServiceAccountToken(saToken AMGServiceAccountToken) error {
	log.Info()
	log.InfoLight("Removing service account token for service account ", saToken.ServiceAccountID)

	_, err := a.Client.DeleteWorkspaceServiceAccountToken(&managedgrafana.DeleteWorkspaceServiceAccountTokenInput{
		ServiceAccountId: aws.String(saToken.ServiceAccountID),
		TokenId:          aws.String(saToken.SATokenID),
		WorkspaceId:      aws.String(saToken.WorkspaceID),
	})

	if err != nil {
		log.Error(err)
	}
	return err
}
