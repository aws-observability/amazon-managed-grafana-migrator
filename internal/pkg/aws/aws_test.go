package aws

import (
	"errors"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws/mocks"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getFakeAPIKey(_ *testing.T) AMGApiKey {
	return AMGApiKey{
		KeyName:     "amg-migrator-", //keyname has a currentime millisecond suffix
		APIKey:      "fakekey",
		WorkspaceID: "g-abcdef1234",
	}
}

func TestAMG_ListWorkspaces(t *testing.T) {

	tests := map[string]struct {
		callMock func(m *mocks.Mockapi)

		expectedWorkspaces int
		expectedError      error
	}{
		"error listing wx": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().ListWorkspaces(&managedgrafana.ListWorkspacesInput{}).Return(
					nil, errors.New("error listing workspaces"),
				)
			},
			expectedWorkspaces: 0,
			expectedError:      errors.New("error listing workspaces"),
		},
		"listing wx": {
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().ListWorkspaces(&managedgrafana.ListWorkspacesInput{}).Return(
					&managedgrafana.ListWorkspacesOutput{
						Workspaces: []*managedgrafana.WorkspaceSummary{
							{
								Id:             aws.String("g-abcdef1234"),
								Name:           aws.String("test"),
								GrafanaVersion: aws.String("8.4"),
								Endpoint:       aws.String("http://g-abcdef1234.us-east-1.grafana.amazonaws.com"),
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
			// GIVEN
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockapi(ctrl)
			tc.callMock(mock)

			client := AMG{
				Client: mock,
			}
			wx, err := client.ListWorkspaces()

			require.Equal(t, tc.expectedWorkspaces, len(wx))
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAMG_CreateWorkspaceApiKey(t *testing.T) {

	tests := map[string]struct {
		workspaceID string
		callMock    func(m *mocks.Mockapi)

		expectedAPIKey AMGApiKey
		expectedError  error
	}{
		"error creating api key": {
			workspaceID: "g-abcdef1234",
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().CreateWorkspaceApiKey(gomock.Any()).Return(nil, errors.New("error creating api key"))
			},
			expectedAPIKey: AMGApiKey{},
			expectedError:  errors.New("error creating api key"),
		},
		"creating api key": {
			workspaceID: "g-abcdef1234",
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().CreateWorkspaceApiKey(gomock.Any()).Return(&managedgrafana.CreateWorkspaceApiKeyOutput{
					Key:         aws.String("fakekey"),
					WorkspaceId: aws.String("g-abcdef1234"),
				}, nil)
			},
			expectedAPIKey: getFakeAPIKey(t),
			expectedError:  nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockapi(ctrl)
			tc.callMock(mock)

			client := AMG{
				Client: mock,
			}
			apiKey, err := client.CreateWorkspaceApiKey(tc.workspaceID)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedAPIKey.APIKey, apiKey.APIKey)
				require.Equal(t, tc.expectedAPIKey.WorkspaceID, apiKey.WorkspaceID)
				require.Contains(t, apiKey.KeyName, tc.expectedAPIKey.KeyName)
			}
		})
	}
}

func TestAMG_DeleteWorkspaceApiKey(t *testing.T) {

	tests := map[string]struct {
		apiKey   AMGApiKey
		callMock func(m *mocks.Mockapi)

		expectedError error
	}{
		"error deleting api key": {
			apiKey: getFakeAPIKey(t),
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().DeleteWorkspaceApiKey(gomock.Any()).Return(nil, errors.New("error deleting api key"))
			},
			expectedError: errors.New("error deleting api key"),
		},
		"deleting api key": {
			apiKey: getFakeAPIKey(t),
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().DeleteWorkspaceApiKey(gomock.Any()).Return(&managedgrafana.DeleteWorkspaceApiKeyOutput{
					KeyName:     aws.String("fakekey"),
					WorkspaceId: aws.String("g-abcdef1234"),
				}, nil)
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

			client := AMG{
				Client: mock,
			}
			err := client.DeleteWorkspaceApiKey(tc.apiKey)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
