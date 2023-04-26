package app

import (
	"errors"
	"testing"

	"amazon-managed-grafana-migrator/internal/pkg/app/mocks"

	"github.com/golang/mock/gomock"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/require"
)

func TestApp_migrateDataSources(t *testing.T) {
	//test cases
	tests := map[string]struct {
		callMockSrc        func(m *mocks.Mockapi)
		callMockDst        func(m *mocks.Mockapi)
		expectedMigratedDs int
		expectedError      error
	}{
		"error getting ds from src": {
			callMockSrc: func(m *mocks.Mockapi) {
				m.EXPECT().DataSources().Return(nil, errors.New("some error")).AnyTimes()
			},
			callMockDst:        func(m *mocks.Mockapi) {},
			expectedMigratedDs: 0,
			expectedError:      errors.New("some error"),
		},
		"syncing ds should continue on error": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockDataSource(t, m)
			},
			callMockDst: func(m *mocks.Mockapi) {
				ds := gapi.DataSource{
					ID:   1,
					UID:  "uid",
					Name: "test-ds",
					URL:  "http://test.com/ds",
				}
				m.EXPECT().NewDataSource(&ds).Return(
					int64(1), errors.New("error while creating ds in dest"),
				).AnyTimes()
			},
			expectedMigratedDs: 0,
			expectedError:      nil,
		},
		"syncing ds": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockDataSource(t, m)
			},
			callMockDst: func(m *mocks.Mockapi) {
				mockNewDatasource(t, m)
			},
			expectedMigratedDs: 1,
			expectedError:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// mocks gapi.Client for src and dst grafana endpoints
			mockSrc := mocks.NewMockapi(ctrl)
			mockDst := mocks.NewMockapi(ctrl)

			tc.callMockSrc(mockSrc)
			tc.callMockDst(mockDst)

			app := App{
				Src: mockSrc,
				Dst: mockDst,
			}

			migratedDs, err := app.migrateDataSources()

			require.Equal(t, migratedDs, tc.expectedMigratedDs)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
