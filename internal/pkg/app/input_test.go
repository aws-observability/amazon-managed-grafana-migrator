package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGrafanaInput(t *testing.T) {
	tests := map[string]struct {
		wkspEndpoint     string
		url              string
		serviceAccountID string
		apiKey           string
		expectedIsAMG    bool
		expectedError    bool
	}{
		"valid AMG workspace": {
			wkspEndpoint:     "g-abcdef1234.grafana-workspace.us-east-1.amazonaws.com",
			serviceAccountID: "sa-1",
			expectedIsAMG:    true,
		},
		"valid external grafana": {
			url:           "https://grafana.example.com",
			apiKey:        "my-api-key",
			expectedIsAMG: false,
		},
		"invalid AMG endpoint": {
			wkspEndpoint:  "invalid-endpoint",
			expectedError: true,
		},
		"missing both inputs": {
			expectedError: true,
		},
		"url without api key": {
			url:           "https://grafana.example.com",
			expectedError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			input, err := NewGrafanaInput(tc.wkspEndpoint, tc.url, tc.serviceAccountID, tc.apiKey)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedIsAMG, input.IsAMG)
			}
		})
	}
}
