package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscover(t *testing.T) {
	tests := map[string]struct {
		input         string
		expectedError string
	}{
		"no region": {
			input:         "",
			expectedError: "missing AWS region",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := discover(tc.input)
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
