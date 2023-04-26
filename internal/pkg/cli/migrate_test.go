package cli

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestGetApiKeyIntegration(t *testing.T) {
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test, set environment variable INTEGRATION=true to enable")
// 	}
// 	//test cases
// 	tests := map[string]struct {
// 		mockSession *session.Session
// 		input       grafanaInput

// 		expectedResult amg.AMGApiKey
// 	}{
// 		"grafana workspace": {

// 			mockSession: session.Must(session.NewSessionWithOptions(session.Options{
// 				SharedConfigState: session.SharedConfigEnable,
// 			})),
// 			input: grafanaInput{
// 				workspaceID: "g-abcdef1234", //get from env
// 				region:      "eu-central-1",
// 				isAMG:       true,
// 			},

// 			//expected result should be
// 			expectedResult: amg.AMGApiKey{
// 				KeyName:     "pelican-",
// 				APIKey:      "fakeAPIKey",
// 				WorkspaceID: "g-abcdef1234",
// 			},
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			apiKey := tc.input.getApiKey(tc.mockSession)
// 			require.Equal(t, tc.expectedResult, apiKey)
// 		})
// 	}

// }

func TestMigrateCmd(t *testing.T) {

	//test cases
	tests := map[string]struct {
		args     []string
		expected string
	}{
		"no args": {
			args:     []string{},
			expected: "migrateUsage:",
		},
		"amg to amg error": {
			args:     []string{"--src", "g-abcdefg123", "--dst", "g-abcdefg234"},
			expected: "migrateUsage:",
		},
		"amg to oss error": {
			args:     []string{"--src", "g-abcdefg234.grafana-workspace.eu-central-1.amazonaws.com", "--dst", "https://grafana.example.com"},
			expected: "migrateUsage:",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			cmd := BuildMigrateCmd()
			cmd.SetArgs(tc.args)
			b := bytes.NewBufferString("migrate")
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
