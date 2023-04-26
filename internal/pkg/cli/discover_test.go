package cli

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscoverCmd(t *testing.T) {

	//test cases
	tests := map[string]struct {
		input    string
		expected string
	}{
		"no args": {
			input:    "discover",
			expected: "discoverUsage:",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			cmd := BuildDiscoverCmd()
			b := bytes.NewBufferString(tc.input)
			cmd.SetOut(b)
			cmd.Execute()
			out, err := io.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}
			require.Contains(t, string(out), tc.expected)
		})
	}
}

func TestDiscover(t *testing.T) {
	//test cases
	tests := map[string]struct {
		input         string
		expectedError error
	}{
		"no args": {
			input:         "",
			expectedError: fmt.Errorf("missing AWS region"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			err := discover(tc.input)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
