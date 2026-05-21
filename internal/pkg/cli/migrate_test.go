package cli

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrateCmd(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected string
	}{
		"no args": {
			args:     []string{},
			expected: "migrateUsage:",
		},
		"invalid src endpoint": {
			args:     []string{"--src", "g-abcdefg123", "--dst", "g-abcdefg234.grafana-workspace.us-east-1.amazonaws.com", "--dst-service-account-id", "sa-1"},
			expected: "migrateUsage:",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
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
