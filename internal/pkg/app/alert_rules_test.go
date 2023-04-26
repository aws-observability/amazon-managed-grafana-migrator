package app

import (
	"amazon-managed-grafana-migrator/internal/pkg/app/mocks"
	"amazon-managed-grafana-migrator/internal/pkg/grafana"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/require"
)

func TestApp_MigrateAlertRules(t *testing.T) {

	tests := map[string]struct {
		callMockDst            func(m *mocks.Mockapi)
		callMockCustomClient   func(m *mocks.MockcustomAPI)
		folders                []gapi.Folder
		expectedMigratedAlerts int
		expectedError          error
	}{
		"migration error": {
			callMockDst: func(m *mocks.Mockapi) {},
			folders:     generateTestFolders(t),
			callMockCustomClient: func(m *mocks.MockcustomAPI) {
				m.EXPECT().AllRuleGroups().Return(
					nil,
					errors.New("error"),
				).AnyTimes()
			},
			expectedMigratedAlerts: 0,
			expectedError:          errors.New("error"),
		},
		"migration success, failing new rule": {
			callMockDst: func(m *mocks.Mockapi) {
				m.EXPECT().NewAlertRule(gomock.Any()).Return("ok", errors.New("error")).AnyTimes()
			},
			folders: generateTestFolders(t),
			callMockCustomClient: func(m *mocks.MockcustomAPI) {
				var rgx grafana.RuleGroupsByFolder
				json.Unmarshal([]byte(mocks.SampleRulesJSON), &rgx)
				m.EXPECT().AllRuleGroups().Return(
					rgx,
					nil,
				).AnyTimes()
			},
			expectedMigratedAlerts: 0, // 3 alerts in mock data
			expectedError:          nil,
		},
		"migration success": {
			callMockDst: func(m *mocks.Mockapi) {
				m.EXPECT().NewAlertRule(gomock.Any()).Return("ok", nil).AnyTimes()
			},
			folders: generateTestFolders(t),
			callMockCustomClient: func(m *mocks.MockcustomAPI) {
				var rgx grafana.RuleGroupsByFolder
				json.Unmarshal([]byte(mocks.SampleRulesJSON), &rgx)
				m.EXPECT().AllRuleGroups().Return(
					rgx,
					nil,
				).AnyTimes()
			},
			expectedMigratedAlerts: 3, // 3 alerts in mock data
			expectedError:          nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// mocks gapi.Client for src and dst grafana endpoints
			mockDst := mocks.NewMockapi(ctrl)

			// mocks grafana.Client for custom grafana client
			mockcustomAPI := mocks.NewMockcustomAPI(ctrl)

			tc.callMockDst(mockDst)
			tc.callMockCustomClient(mockcustomAPI)

			app := App{
				Src: mockDst,
				Dst: mockDst,
			}

			migratedDs, err := app.migrateAlertRules(tc.folders, CustomGrafanaClient{Client: mockcustomAPI})

			require.Equal(t, tc.expectedMigratedAlerts, migratedDs)
			if tc.expectedError != nil {
				require.EqualError(t, tc.expectedError, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}

}

// func TestApp_migrateAlerts(t *testing.T) {
// 	//test cases
// 	tests := map[string]struct {
// 		callMockSrc            func(m *mocks.Mockapi)
// 		callMockDst            func(m *mocks.Mockapi)
// 		expectedMigratedAlerts int
// 		expectedError          error
// 	}{
// 		"error getting dashboards from src": {
// 			callMockSrc: func(m *mocks.Mockapi) {
// 				m.EXPECT().Dashboards().Return(
// 					nil,
// 					errors.New("some error"),
// 				).AnyTimes()
// 			},
// 			callMockDst: func(m *mocks.Mockapi) {
// 				m.EXPECT().Dashboards().Return(
// 					nil,
// 					nil,
// 				).AnyTimes()
// 			},
// 			expectedMigratedAlerts: 0,
// 			expectedError:          errors.New("some error"),
// 		},
// 		"syncing ds should continue on error": {
// 			callMockSrc: func(m *mocks.Mockapi) {
// 				dsResponse := gapi.FolderDashboardSearchResponse{
// 					FolderID:  1,
// 					FolderUID: "uid",
// 					Title:     "Test Dashboard",
// 					UID:       "uid",
// 				}
// 				ds := gapi.Dashboard{
// 					FolderID:  1,
// 					FolderUID: "uid",
// 					Overwrite: true,
// 					Meta:      gapi.DashboardMeta{},
// 				}
// 				m.EXPECT().Dashboards().Return([]gapi.FolderDashboardSearchResponse{dsResponse}, nil).AnyTimes()
// 				m.EXPECT().DashboardByUID("uid").Return(&ds, errors.New("some error")).AnyTimes()
// 			},
// 			callMockDst:            func(m *mocks.Mockapi) {},
// 			expectedMigratedAlerts: 0,
// 			expectedError:          nil,
// 		},
// 		// TODO: test this in integration, as the grafana behaviour cannot be mocked
// 		// "syncing ds new dashboard error": {
// 		// 	callMockSrc: func(m *mocks.Mockapi) {
// 		// 		dsResponse := gapi.FolderDashboardSearchResponse{
// 		// 			FolderID:  1,
// 		// 			FolderUID: "uid",
// 		// 			Title:     "Test Dashboard",
// 		// 			UID:       "uid",
// 		// 		}
// 		// 		ds := gapi.Dashboard{
// 		// 			FolderID:  1,
// 		// 			FolderUID: "uid",
// 		// 			Overwrite: true,
// 		// 			Meta:      gapi.DashboardMeta{},
// 		// 		}
// 		// 		m.EXPECT().Dashboards().Return([]gapi.FolderDashboardSearchResponse{dsResponse}, nil).AnyTimes()
// 		// 		m.EXPECT().DashboardByUID("uid").Return(&ds, nil).AnyTimes()
// 		// 	},
// 		// 	callMockDst: func(m *mocks.Mockapi) {
// 		// 		ds := gapi.Dashboard{
// 		// 			FolderID:  1,
// 		// 			FolderUID: "uid",
// 		// 			Overwrite: true,
// 		// 			Meta: gapi.DashboardMeta{
// 		// 				Folder: 1,
// 		// 			},
// 		// 		}
// 		// 		m.EXPECT().NewDashboard(&ds).Return(&gapi.DashboardSaveResponse{}, errors.New("Error creating new dashboard")).AnyTimes()
// 		// 	},
// 		// 	expectedMigratedDs: 0,
// 		// 	expectedError:      nil,
// 		// },
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			// mocks gapi.Client for src and dst grafana endpoints
// 			mockSrc := mocks.NewMockapi(ctrl)
// 			mockDst := mocks.NewMockapi(ctrl)

// 			tc.callMockSrc(mockSrc)
// 			tc.callMockDst(mockDst)

// 			app := App{
// 				Src: mockSrc,
// 				Dst: mockDst,
// 			}

// 			migratedDs, err := app.migrateAlerts()

// 			require.Equal(t, migratedDs, tc.expectedMigratedAlerts)
// 			if tc.expectedError != nil {
// 				require.EqualError(t, err, tc.expectedError.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }
