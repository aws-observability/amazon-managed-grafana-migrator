// Package grafana is an http client for api calls authenticated with
// not supported by github.com/grafana/grafana-api-golang-client
// TODO: make a PR to gapi and drop this duplicate implementation
package grafana

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	tests := map[string]struct {
		url           string
		apiKey        string
		expectedError error
	}{
		"malformed url": {
			url:    " http://example.org",
			apiKey: "somekey",

			expectedError: errors.New("http://example.org"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			client, err := New(tc.url, tc.apiKey)

			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.IsType(t, &Client{}, client)
			}
		})
	}
}
