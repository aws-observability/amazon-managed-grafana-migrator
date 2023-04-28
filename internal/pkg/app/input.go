package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

// GrafanaInput holds the infos about the grafana server from the CLI
type GrafanaInput struct {
	URL         string
	WorkspaceID string
	APIKey      string
	Region      string
	IsAMG       bool
	//TODO: remove
	IsGamma bool
}

// GrafanaHTTPClient contains the grafana client and AWS API key
type GrafanaHTTPClient struct {
	Client *gapi.Client
	Key    aws.AMGApiKey
	Input  *GrafanaInput
}

// NewGrafanaInput validate input from command line to return a GrafanaInput object
func NewGrafanaInput(wkspEndpoint, url, apiKey string) (GrafanaInput, error) {

	if wkspEndpoint != "" {
		sx := strings.Split(wkspEndpoint, ".")
		if len(sx) != 5 {
			return GrafanaInput{}, fmt.Errorf("invalid input: workspace should be its DNS endpoint")
		}
		return GrafanaInput{
			WorkspaceID: sx[0],
			Region:      sx[2],
			URL:         wkspEndpoint,
			IsAMG:       true,
			IsGamma:     strings.Contains(sx[1], "gamma"),
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

// getAPIKey create api keys only when provided with a managed grafana ID
func (input *GrafanaInput) getAPIKey(awsgrafanacli *aws.AMG) (aws.AMGApiKey, error) {

	if !input.IsAMG {
		log.Debug("Skipping API key creation for ", input.URL)
		return aws.AMGApiKey{
			APIKey: input.APIKey,
		}, nil
	}

	key, err := awsgrafanacli.CreateWorkspaceApiKey(input.WorkspaceID)
	if err != nil {
		return key, err
	}
	return key, nil
}

// CreateGrafanaAPIClient create a grafana HTTP API client from the input
func (input *GrafanaInput) CreateGrafanaAPIClient(awsgrafanacli *aws.AMG) (*GrafanaHTTPClient, error) {
	var url string

	if input.IsAMG {
		// could be replaced by describe workspace
		url = fmt.Sprintf("https://%s", input.URL)
	} else {
		url = input.URL
	}

	// get API keys
	apiKey, err := input.getAPIKey(awsgrafanacli)
	if err != nil {
		return nil, err
	}

	client, err := gapi.New(url, gapi.Config{APIKey: apiKey.APIKey})
	if err != nil {
		return nil, err
	}
	return &GrafanaHTTPClient{
		Client: client,
		Key:    apiKey,
		Input:  input,
	}, nil
}

// DeleteAPIKeys delete the temporary API key from the AWS grafana workspace
func (input *GrafanaInput) DeleteAPIKeys(awsgrafanacli *aws.AMG, apiKey aws.AMGApiKey) error {
	if !input.IsAMG {
		return nil
	}
	return awsgrafanacli.DeleteWorkspaceApiKey(apiKey)
}
