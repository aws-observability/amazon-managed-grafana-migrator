// Package aws provides a wrapper around the AWS Managed Grafana API
package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
)

type api interface {
	ListWorkspaces(ctx context.Context, input *grafana.ListWorkspacesInput, optFns ...func(*grafana.Options)) (*grafana.ListWorkspacesOutput, error)
	CreateWorkspaceServiceAccountToken(ctx context.Context, input *grafana.CreateWorkspaceServiceAccountTokenInput, optFns ...func(*grafana.Options)) (*grafana.CreateWorkspaceServiceAccountTokenOutput, error)
	DeleteWorkspaceServiceAccountToken(ctx context.Context, input *grafana.DeleteWorkspaceServiceAccountTokenInput, optFns ...func(*grafana.Options)) (*grafana.DeleteWorkspaceServiceAccountTokenOutput, error)
	DescribeWorkspace(ctx context.Context, input *grafana.DescribeWorkspaceInput, optFns ...func(*grafana.Options)) (*grafana.DescribeWorkspaceOutput, error)
}

// GrafanaAuth provides the authentication token for a Grafana workspace
type GrafanaAuth interface {
	GetAuth() string
}

// AMGServiceAccountToken is a GrafanaAuth struct for service account tokens
type AMGServiceAccountToken struct {
	ServiceAccountID, SATokenID, Token, WorkspaceID string
}

// GetAuth returns the SA token for a workspace
func (t AMGServiceAccountToken) GetAuth() string {
	return t.Token
}

// ExternalAPIKey is a GrafanaAuth for non-AMG Grafana servers using a static API key
type ExternalAPIKey struct {
	APIKey string
}

// GetAuth returns the API key
func (k ExternalAPIKey) GetAuth() string {
	return k.APIKey
}

// AMG is an AWS SDK client for AMG APIs
type AMG struct {
	Client api
}

// Workspace contains information about a Grafana workspace
type Workspace struct {
	ID       string
	Name     string
	Version  string
	Endpoint string
	Status   types.WorkspaceStatus
}

// New creates a new AMG client from an aws.Config
func New(cfg aws.Config, region string) *AMG {
	return &AMG{
		Client: grafana.NewFromConfig(cfg, func(o *grafana.Options) {
			o.Region = region
		}),
	}
}

// ListWorkspaces lists all workspaces in the region
func (a *AMG) ListWorkspaces(ctx context.Context) ([]Workspace, error) {
	response, err := a.Client.ListWorkspaces(ctx, &grafana.ListWorkspacesInput{})
	if err != nil {
		return nil, err
	}

	wx := make([]Workspace, 0, len(response.Workspaces))
	for _, ws := range response.Workspaces {
		wx = append(wx, Workspace{
			ID:       aws.ToString(ws.Id),
			Name:     aws.ToString(ws.Name),
			Version:  aws.ToString(ws.GrafanaVersion),
			Endpoint: aws.ToString(ws.Endpoint),
			Status:   ws.Status,
		})
	}
	return wx, nil
}

// DescribeWorkspace returns information about a workspace
func (a *AMG) DescribeWorkspace(ctx context.Context, workspaceID string) (Workspace, error) {
	res, err := a.Client.DescribeWorkspace(ctx, &grafana.DescribeWorkspaceInput{
		WorkspaceId: aws.String(workspaceID),
	})
	if err != nil {
		return Workspace{}, err
	}
	return Workspace{
		ID:       workspaceID,
		Name:     aws.ToString(res.Workspace.Name),
		Version:  aws.ToString(res.Workspace.GrafanaVersion),
		Endpoint: aws.ToString(res.Workspace.Endpoint),
		Status:   res.Workspace.Status,
	}, nil
}

// CreateServiceAccountToken creates a temporary service account token
func (a *AMG) CreateServiceAccountToken(ctx context.Context, workspaceID, serviceAccountID string) (AMGServiceAccountToken, error) {
	log.InfoLight("Creating service account token for service account ", serviceAccountID)

	tokenName := fmt.Sprintf("amg-migrator-%d", time.Now().UTC().UnixMilli())
	ttl := int32(30 * 60) // 30 minutes

	resp, err := a.Client.CreateWorkspaceServiceAccountToken(ctx, &grafana.CreateWorkspaceServiceAccountTokenInput{
		Name:             aws.String(tokenName),
		SecondsToLive:    aws.Int32(ttl),
		WorkspaceId:      aws.String(workspaceID),
		ServiceAccountId: aws.String(serviceAccountID),
	})
	if err != nil {
		return AMGServiceAccountToken{}, err
	}

	return AMGServiceAccountToken{
		ServiceAccountID: serviceAccountID,
		SATokenID:        aws.ToString(resp.ServiceAccountToken.Id),
		Token:            aws.ToString(resp.ServiceAccountToken.Key),
		WorkspaceID:      workspaceID,
	}, nil
}

// DeleteServiceAccountToken removes a service account token
func (a *AMG) DeleteServiceAccountToken(ctx context.Context, saToken AMGServiceAccountToken) error {
	log.InfoLight("Removing service account token for service account ", saToken.ServiceAccountID)

	_, err := a.Client.DeleteWorkspaceServiceAccountToken(ctx, &grafana.DeleteWorkspaceServiceAccountTokenInput{
		ServiceAccountId: aws.String(saToken.ServiceAccountID),
		TokenId:          aws.String(saToken.SATokenID),
		WorkspaceId:      aws.String(saToken.WorkspaceID),
	})
	if err != nil {
		log.Error(err)
	}
	return err
}
