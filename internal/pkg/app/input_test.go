package app

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/aws/mocks"

	awssdk "github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/golang/mock/gomock"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/require"
)

func TestInput_NewGrafanaInput(t *testing.T) {
	//test cases
	tests := map[string]struct {
		wkspEndpoint string
		url          string
		apiKey       string

		expectedResult GrafanaInput
		expectedError  error
	}{
		"grafana workspace ID": {
			// this should be an endpoint instead
			wkspEndpoint: "g-abcdefg234",
			// workspace ID is mutually exclusive with URL/API key
			url:    "",
			apiKey: "",

			//expected input should be
			expectedResult: GrafanaInput{},
			//expected error should be
			expectedError: fmt.Errorf("invalid input: workspace should be its DNS endpoint"),
		},
		"grafana workspace": {
			wkspEndpoint: "g-abcdefg234.grafana-workspace.eu-central-1.amazonaws.com",
			// workspace ID is mutually exclusive with URL/API key
			url:    "",
			apiKey: "",

			//expected input should be
			expectedResult: GrafanaInput{
				WorkspaceID: "g-abcdefg234",
				Region:      "eu-central-1",
				IsAMG:       true,
				URL:         "g-abcdefg234.grafana-workspace.eu-central-1.amazonaws.com",
			},
			//expected error should be
			expectedError: nil,
		},
		"grafana oss": {
			wkspEndpoint: "",
			// workspace ID is mutually exclusive with URL/API key
			url:    "https://grafana.example.com",
			apiKey: "fakeAPIKey",

			//expected input should be
			expectedResult: GrafanaInput{
				URL:    "https://grafana.example.com",
				APIKey: "fakeAPIKey",
				IsAMG:  false,
			},
			//expected error should be
			expectedError: nil,
		},
		"no input": {
			wkspEndpoint: "",
			url:          "",
			apiKey:       "",

			//expected input should be
			expectedResult: GrafanaInput{},
			//expected error should be
			expectedError: errors.New("invalid input"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			input, err := NewGrafanaInput(tc.wkspEndpoint, tc.url, tc.apiKey)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, input)
			}
		})
	}

}

func generateGrafanaInput(_ *testing.T, isAMG bool) GrafanaInput {
	if isAMG {
		return GrafanaInput{
			WorkspaceID: "g-abcdefg234",
			Region:      "eu-central-1",
			IsAMG:       true,
			APIKey:      "fakeAPIKey",
			URL:         "g-abcdefg234.grafana-workspace.eu-central-1.amazonaws.com",
		}
	}
	return GrafanaInput{
		WorkspaceID: "",
		Region:      "",
		IsAMG:       false,
		APIKey:      "fakeAPIKey",
		URL:         "https://grafana.example.com",
	}
}

func TestInput_GetApiKey(t *testing.T) {
	//test cases
	tests := map[string]struct {
		mockSession *session.Session
		input       GrafanaInput

		expectedResult aws.AMGApiKey
		expectedError  error
	}{
		"grafana oss": {
			mockSession: nil,
			input:       generateGrafanaInput(t, false),
			//expected error should be
			expectedResult: aws.AMGApiKey{
				APIKey: "fakeAPIKey",
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockapi(ctrl)
			// tc.callMock(mock)

			awsgrafanacli := aws.AMG{
				Client: mock,
			}

			apiKey, err := tc.input.getAPIKey(&awsgrafanacli)
			require.Equal(t, tc.expectedResult, apiKey)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestInput_CreateGrafanaAPIClient(t *testing.T) {

	tests := map[string]struct {
		input    GrafanaInput
		callMock func(m *mocks.Mockapi)

		expectedError error
	}{
		"oss client": {
			input: generateGrafanaInput(t, false),
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().CreateWorkspaceApiKey(gomock.Any()).Return(nil, errors.New("error creating api key")).AnyTimes()
			},
			expectedError: nil,
		},
		"amg client": {
			input: generateGrafanaInput(t, true),
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().CreateWorkspaceApiKey(gomock.Any()).Return(&managedgrafana.CreateWorkspaceApiKeyOutput{
					Key:         awssdk.String("fakekey"),
					WorkspaceId: awssdk.String("g-abcdef1234"),
				}, nil).AnyTimes()
				m.EXPECT().DeleteWorkspaceApiKey(gomock.Any()).Return(&managedgrafana.DeleteWorkspaceApiKeyOutput{
					KeyName:     awssdk.String("amg-migrator-"),
					WorkspaceId: awssdk.String("g-abcdef1234"),
				}, nil).AnyTimes()
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

			awsgrafanacli := aws.AMG{
				Client: mock,
			}

			client, err := tc.input.CreateGrafanaAPIClient(&awsgrafanacli)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.IsType(t, &gapi.Client{}, client.Client)
			}
		})
	}

}

func TestInput_DeleteAPIKeys(t *testing.T) {

	tests := map[string]struct {
		input     GrafanaInput
		amgAPIKey aws.AMGApiKey
		callMock  func(m *mocks.Mockapi)

		expectedError error
	}{
		"amg workspace": {
			input: generateGrafanaInput(t, true),
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().DeleteWorkspaceApiKey(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			expectedError: nil,
		},
		"not an amg workspace": {
			input: generateGrafanaInput(t, false),
			callMock: func(m *mocks.Mockapi) {
				m.EXPECT().DeleteWorkspaceApiKey(gomock.Any()).Return(nil, nil).AnyTimes()
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

			awsgrafanacli := aws.AMG{
				Client: mock,
			}

			err := tc.input.DeleteAPIKeys(&awsgrafanacli, tc.amgAPIKey)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}

}
