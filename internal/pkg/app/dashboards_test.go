package app

import (
	"errors"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/app/mocks"

	"github.com/golang/mock/gomock"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/require"
)

func TestApp_migrateDashboards(t *testing.T) {
	//test cases
	tests := map[string]struct {
		callMockSrc        func(m *mocks.Mockapi)
		callMockDst        func(m *mocks.Mockapi)
		expectedMigratedDs int
		expectedError      error
	}{
		"error getting dashboards from src": {
			callMockSrc: func(m *mocks.Mockapi) {
				m.EXPECT().Dashboards().Return(
					nil,
					errors.New("some error"),
				).AnyTimes()
			},
			callMockDst: func(m *mocks.Mockapi) {
				m.EXPECT().Dashboards().Return(
					nil,
					nil,
				).AnyTimes()
			},
			expectedMigratedDs: 0,
			expectedError:      errors.New("some error"),
		},
		"syncing ds should continue on error": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockDashboards(t, m)
			},
			callMockDst:        func(m *mocks.Mockapi) {},
			expectedMigratedDs: 0,
			expectedError:      nil,
		},
		// TODO: test this in integration, as the grafana behaviour cannot be mocked
		// "syncing ds new dashboard error": {
		// 	callMockSrc: func(m *mocks.Mockapi) {
		// 		dsResponse := gapi.FolderDashboardSearchResponse{
		// 			FolderID:  1,
		// 			FolderUID: "uid",
		// 			Title:     "Test Dashboard",
		// 			UID:       "uid",
		// 		}
		// 		ds := gapi.Dashboard{
		// 			FolderID:  1,
		// 			FolderUID: "uid",
		// 			Overwrite: true,
		// 			Meta:      gapi.DashboardMeta{},
		// 		}
		// 		m.EXPECT().Dashboards().Return([]gapi.FolderDashboardSearchResponse{dsResponse}, nil).AnyTimes()
		// 		m.EXPECT().DashboardByUID("uid").Return(&ds, nil).AnyTimes()
		// 	},
		// 	callMockDst: func(m *mocks.Mockapi) {
		// 		ds := gapi.Dashboard{
		// 			FolderID:  1,
		// 			FolderUID: "uid",
		// 			Overwrite: true,
		// 			Meta: gapi.DashboardMeta{
		// 				Folder: 1,
		// 			},
		// 		}
		// 		m.EXPECT().NewDashboard(&ds).Return(&gapi.DashboardSaveResponse{}, errors.New("Error creating new dashboard")).AnyTimes()
		// 	},
		// 	expectedMigratedDs: 0,
		// 	expectedError:      nil,
		// },
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

			f := gapi.Folder{
				ID:    1,
				UID:   "uid",
				Title: "test",
				URL:   "http://test.com",
			}

			migratedDs, err := app.migrateDashboards(&[]gapi.Folder{f})

			require.Equal(t, migratedDs, tc.expectedMigratedDs)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
