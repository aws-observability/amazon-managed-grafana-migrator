// Package grafana is an http client for api calls authenticated with
// not supported by github.com/grafana/grafana-api-golang-client
// TODO: make a PR to gapi and drop this duplicate implementation
package grafana

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/hashicorp/go-cleanhttp"
)

// HTTPClient is an interface to go http package
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a Grafana API client.
type Client struct {
	APIKey  string
	baseURL url.URL
	client  HTTPClient
}

// New creates a new Grafana client.
func New(baseURL, apiKey string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		APIKey:  apiKey,
		baseURL: *u,
		client:  cleanhttp.DefaultClient(),
	}, nil
}

// Get sends a get requets and unmarshall result into responseStruct
func (c *Client) get(requestPath string, responseStruct interface{}) error {
	url := c.baseURL
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	url.Path = path.Join(url.Path, requestPath)

	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read the body (even on non-successful HTTP status codes), as that's what the unit tests expect
	bodyContents, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyContents, responseStruct)
	if err != nil {
		return err
	}

	return nil
}
