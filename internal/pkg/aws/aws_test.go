package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws/mocks"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAMG_ListWorkspaces(t *testing.T) {
	tests := map[string]struct {
		callMock           func(m *mocks.Mockapi)
		expectedWorkspaces int
		expectedError      error
	}{
		"error listing workspaces": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().ListWorkspaces(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("error listing workspaces"),
				)
			},
			expectedWorkspaces: 0,
			expectedError:      errors.New("error listing workspaces"),
		},
		"listing workspaces": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().ListWorkspaces(gomock.Any(), gomock.Any()).Return(
					&grafana.ListWorkspacesOutput{
						Workspaces: []types.WorkspaceSummary{
							{
								Id:             aws.String("g-abcdef1234"),
								Name:           aws.String("test"),
								GrafanaVersion: aws.String("10.4"),
								Endpoint:       aws.String("g-abcdef1234.grafana-workspace.us-east-1.amazonaws.com"),
								Status:         types.WorkspaceStatusActive,
							},
						},
					}, nil,
				)
			},
			expectedWorkspaces: 1,
			expectedError:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockapi(ctrl)
			tc.callMock(mock)

			client := AMG{Client: mock}
			wx, err := client.ListWorkspaces(context.Background())

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Len(t, wx, tc.expectedWorkspaces)
			}
		})
	}
}

func TestAMG_CreateServiceAccountToken(t *testing.T) {
	tests := map[string]struct {
		callMock      func(m *mocks.Mockapi)
		expectedError error
	}{
		"error creating token": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().CreateWorkspaceServiceAccountToken(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("error creating token"),
				)
			},
			expectedError: errors.New("error creating token"),
		},
		"creating token": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().CreateWorkspaceServiceAccountToken(gomock.Any(), gomock.Any()).Return(
					&grafana.CreateWorkspaceServiceAccountTokenOutput{
						ServiceAccountToken: &types.ServiceAccountTokenSummaryWithKey{
							Id:  aws.String("token-123"),
							Key: aws.String("fake-token-key"),
						},
					}, nil,
				)
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockapi(ctrl)
			tc.callMock(mock)

			client := AMG{Client: mock}
			token, err := client.CreateServiceAccountToken(context.Background(), "g-abcdef1234", "sa-1")

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, "fake-token-key", token.Token)
				require.Equal(t, "token-123", token.SATokenID)
				require.Equal(t, "sa-1", token.ServiceAccountID)
			}
		})
	}
}

func TestAMG_DeleteServiceAccountToken(t *testing.T) {
	tests := map[string]struct {
		callMock      func(m *mocks.Mockapi)
		expectedError error
	}{
		"error deleting token": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().DeleteWorkspaceServiceAccountToken(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("error deleting token"),
				)
			},
			expectedError: errors.New("error deleting token"),
		},
		"deleting token": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().DeleteWorkspaceServiceAccountToken(gomock.Any(), gomock.Any()).Return(
					&grafana.DeleteWorkspaceServiceAccountTokenOutput{}, nil,
				)
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockapi(ctrl)
			tc.callMock(mock)

			client := AMG{Client: mock}
			saToken := AMGServiceAccountToken{
				ServiceAccountID: "sa-1",
				SATokenID:        "token-123",
				Token:            "fake-token-key",
				WorkspaceID:      "g-abcdef1234",
			}
			err := client.DeleteServiceAccountToken(context.Background(), saToken)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
