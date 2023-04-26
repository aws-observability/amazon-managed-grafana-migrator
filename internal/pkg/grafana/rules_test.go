package grafana

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/grafana/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRules_AllRuleGroups(t *testing.T) {
	tests := map[string]struct {
		callMock           func(m *mocks.MockHTTPClient)
		expectedRuleGroups int
		expectedError      error
	}{
		"success": {
			callMock: func(m *mocks.MockHTTPClient) {
				json := `{"Folder": [{"name": "rulegroup", "interval": "1m"}]}`
				// create a new reader with that JSON
				r := io.NopCloser(bytes.NewReader([]byte(json)))
				m.EXPECT().Do(gomock.Any()).Return(
					&http.Response{
						StatusCode: 200,
						Body:       r,
					}, nil,
				).AnyTimes()
			},
			expectedRuleGroups: 1,
			expectedError:      nil,
		},
		"unmarshall error": {
			callMock: func(m *mocks.MockHTTPClient) {
				json := `{"Hello": "World"}`
				// create a new reader with that JSON
				r := io.NopCloser(bytes.NewReader([]byte(json)))
				m.EXPECT().Do(gomock.Any()).Return(
					&http.Response{
						StatusCode: 200,
						Body:       r,
					}, nil,
				).AnyTimes()
			},
			expectedRuleGroups: 0,
			expectedError:      errors.New("json: cannot unmarshal string into Go value of type []grafana.RuleGroup"),
		},
		"http error": {
			callMock: func(m *mocks.MockHTTPClient) {
				m.EXPECT().Do(gomock.Any()).Return(
					&http.Response{
						StatusCode: 400,
						Body:       nil,
					}, errors.New("http error"),
				).AnyTimes()
			},
			expectedRuleGroups: 0,
			expectedError:      errors.New("http error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockHTTPClient(ctrl)
			tc.callMock(mock)

			client, _ := New("localhost", "apikey")
			client.client = mock

			rgx, err := client.AllRuleGroups()

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Len(t, rgx, tc.expectedRuleGroups)
			}
		})
	}
}
