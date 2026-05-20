package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"
	"github.com/hashicorp/go-cleanhttp"
)

const (
	AMG_V12 = "12.4"
	AMG_V10 = "10.4"
	AMG_V9  = "9.4"
)

// GrafanaInput holds the infos about the grafana server from the CLI
type GrafanaInput struct {
	URL              string
	WorkspaceID      string
	APIKey           string
	Region           string
	ServiceAccountID string
	WorkspaceVersion string
	IsAMG            bool
}

// GrafanaHTTPClient contains the grafana HTTP client and AWS auth token
type GrafanaHTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Auth       aws.GrafanaAuth
	Input      *GrafanaInput
}

// NewGrafanaInput validates input from command line to return a GrafanaInput object
func NewGrafanaInput(wkspEndpoint, url, serviceAccountID, apiKey string) (GrafanaInput, error) {
	if wkspEndpoint != "" {
		sx := strings.Split(wkspEndpoint, ".")
		if len(sx) != 5 {
			return GrafanaInput{}, fmt.Errorf("invalid input: workspace should be its DNS endpoint")
		}
		return GrafanaInput{
			WorkspaceID:      sx[0],
			Region:           sx[2],
			URL:              wkspEndpoint,
			ServiceAccountID: serviceAccountID,
			IsAMG:            true,
		}, nil
	} else if url != "" && apiKey != "" {
		return GrafanaInput{
			URL:    url,
			APIKey: apiKey,
			IsAMG:  false,
		}, nil
	}

	return GrafanaInput{}, errors.New("invalid input")
}

// getGrafanaAuthToken creates a service account token for AMG workspaces
func (input *GrafanaInput) getGrafanaAuthToken(ctx context.Context, awsgrafanacli *aws.AMG) (aws.GrafanaAuth, error) {
	if !input.IsAMG {
		log.InfoLight("Using provided API key for ", input.URL)
		return aws.ExternalAPIKey{APIKey: input.APIKey}, nil
	}

	wksp, err := awsgrafanacli.DescribeWorkspace(ctx, input.WorkspaceID)
	if err == nil {
		input.WorkspaceVersion = wksp.Version
	}

	if input.ServiceAccountID == "" {
		return nil, errors.New("input error: service account ID is required for AMG workspaces (v9+), run migrate -h for help")
	}

	return awsgrafanacli.CreateServiceAccountToken(ctx, input.WorkspaceID, input.ServiceAccountID)
}

// CreateGrafanaHTTPClient creates a Grafana HTTP client from the input
func (input *GrafanaInput) CreateGrafanaHTTPClient(ctx context.Context, awsgrafanacli *aws.AMG) (*GrafanaHTTPClient, error) {
	var baseURL string
	if input.IsAMG {
		baseURL = fmt.Sprintf("https://%s", input.URL)
	} else {
		baseURL = input.URL
	}

	auth, err := input.getGrafanaAuthToken(ctx, awsgrafanacli)
	if err != nil {
		return nil, err
	}

	return &GrafanaHTTPClient{
		BaseURL:    baseURL,
		HTTPClient: cleanhttp.DefaultClient(),
		Auth:       auth,
		Input:      input,
	}, nil
}

// DeleteGrafanaAuth deletes the temporary service account token
func (input *GrafanaInput) DeleteGrafanaAuth(ctx context.Context, awsgrafanacli *aws.AMG, auth aws.GrafanaAuth) error {
	if !input.IsAMG {
		return nil
	}
	return awsgrafanacli.DeleteServiceAccountToken(ctx, auth.(aws.AMGServiceAccountToken))
}
