package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

const (
	AMG_V10 = "10.4"
	AMG_V9  = "9.4"
	AMG_V8  = "8.4"
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
	IsGamma          bool
}

// GrafanaHTTPClient contains the grafana client and AWS API key
type GrafanaHTTPClient struct {
	Client *gapi.Client
	Auth   aws.GrafanaAuth
	Input  *GrafanaInput
}

// NewGrafanaInput validate input from command line to return a GrafanaInput object
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
			IsGamma:          strings.Contains(sx[1], "gamma"),
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

// getGrafanaAuthToken create Grafana api keys only when provided with a managed grafana ID
// if service account is provided, it will create a service account token
func (input *GrafanaInput) getGrafanaAuthToken(awsgrafanacli *aws.AMG) (aws.GrafanaAuth, error) {

	if !input.IsAMG {
		log.InfoLight("Skipping API key creation for ", input.URL)
		return aws.AMGApiKey{
			APIKey: input.APIKey,
		}, nil
	}

	wksp, err := awsgrafanacli.DescribeWorkspace(input.WorkspaceID)
	if err == nil {
		input.WorkspaceVersion = wksp.Version
	}

	// forcing V10 to use service account token
	if input.WorkspaceVersion == AMG_V10 && input.ServiceAccountID == "" {
		return nil, errors.New("input error: service account ID is required for AMG v10, run migrate -h for help")
	}

	// creating service account token if service account is provided
	if input.ServiceAccountID != "" {
		return awsgrafanacli.CreateServiceAccountToken(input.WorkspaceID, input.ServiceAccountID)
	}

	// creating temporary API key if no service account is provided
	return awsgrafanacli.CreateWorkspaceApiKey(input.WorkspaceID)
}

// CreateGrafanaAPIClient create a grafana HTTP API client from the input
func (input *GrafanaInput) CreateGrafanaAPIClient(awsgrafanacli *aws.AMG) (*GrafanaHTTPClient, error) {
	var url string

	if input.IsAMG {
		url = fmt.Sprintf("https://%s", input.URL)
	} else {
		url = input.URL
	}

	// get final auth key or token
	apiKey, err := input.getGrafanaAuthToken(awsgrafanacli)
	if err != nil {
		return nil, err
	}

	client, err := gapi.New(url, gapi.Config{APIKey: apiKey.GetAuth()})
	if err != nil {
		return nil, err
	}
	return &GrafanaHTTPClient{
		Client: client,
		Auth:   apiKey,
		Input:  input,
	}, nil
}

// DeleteGrafanaAuth delete the temporary API key from the AWS grafana workspace
func (input *GrafanaInput) DeleteGrafanaAuth(awsgrafanacli *aws.AMG, auth aws.GrafanaAuth) error {
	if !input.IsAMG {
		return nil
	}

	if input.ServiceAccountID != "" {
		return awsgrafanacli.DeleteServiceAccountToken(auth.(aws.AMGServiceAccountToken))
	}

	return awsgrafanacli.DeleteWorkspaceApiKey(auth.(aws.AMGApiKey))
}
