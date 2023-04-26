package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/app/mocks"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	"github.com/golang/mock/gomock"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/require"
)

func mockDataSource(_ *testing.T, m *mocks.Mockapi) {
	ds := gapi.DataSource{
		ID:   1,
		UID:  "uid",
		Name: "test-ds",
		URL:  "http://test.com/ds",
	}
	m.EXPECT().DataSources().Return([]*gapi.DataSource{&ds}, nil).AnyTimes()
}

func mockNewDatasource(_ *testing.T, m *mocks.Mockapi) {
	ds := gapi.DataSource{
		ID:   1,
		UID:  "uid",
		Name: "test-ds",
		URL:  "http://test.com/ds",
	}
	m.EXPECT().NewDataSource(&ds).Return(
		int64(1), nil,
	).AnyTimes()
}

func mockFolder(_ *testing.T, m *mocks.Mockapi) {
	f := gapi.Folder{
		ID:    1,
		UID:   "uid",
		Title: "test",
		URL:   "http://test.com",
	}
	m.EXPECT().Folders().Return([]gapi.Folder{f}, nil).AnyTimes()
}

func mockNewFolder(_ *testing.T, m *mocks.Mockapi) {
	f := gapi.Folder{
		ID:    1,
		UID:   "uid",
		Title: "test",
		URL:   "http://test.com",
	}
	m.EXPECT().NewFolder(f.Title, f.UID).Return(
		gapi.Folder{}, errors.New("some error while creating folder in dest"),
	)
}
func mockDashboards(_ *testing.T, m *mocks.Mockapi) {
	dsResponse := gapi.FolderDashboardSearchResponse{
		FolderID:  1,
		FolderUID: "uid",
		Title:     "Test Dashboard",
		UID:       "uid",
	}
	ds := gapi.Dashboard{
		FolderID:  1,
		FolderUID: "uid",
		Overwrite: true,
		Meta:      gapi.DashboardMeta{},
	}
	m.EXPECT().Dashboards().Return([]gapi.FolderDashboardSearchResponse{dsResponse}, nil).AnyTimes()
	m.EXPECT().DashboardByUID("uid").Return(&ds, errors.New("some error")).AnyTimes()
}

func TestApp_Run(t *testing.T) {

	tests := map[string]struct {
		callMockSrc          func(m *mocks.Mockapi)
		callMockDst          func(m *mocks.Mockapi)
		callMockCustomClient func(m *mocks.MockcustomAPI)
		expectedError        error
	}{
		"error migrating ds": {
			callMockSrc: func(m *mocks.Mockapi) {
				m.EXPECT().DataSources().Return(nil, errors.New("some error")).AnyTimes()
				m.EXPECT().Folders().Return(nil, errors.New("some error")).AnyTimes()
				m.EXPECT().Dashboards().Return(nil, errors.New("some error")).AnyTimes()
			},
			callMockDst:          func(m *mocks.Mockapi) {},
			callMockCustomClient: func(m *mocks.MockcustomAPI) {},
			expectedError:        errors.New("some error"),
		},
		"migrating ds": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockDataSource(t, m)
				m.EXPECT().Folders().Return(nil, errors.New("some error")).AnyTimes()
				m.EXPECT().Dashboards().Return(nil, errors.New("some error")).AnyTimes()
			},
			callMockDst: func(m *mocks.Mockapi) {
				mockNewDatasource(t, m)
			},
			callMockCustomClient: func(m *mocks.MockcustomAPI) {},
			expectedError:        errors.New("some error"),
		},
		"migrating folders": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockDataSource(t, m)
				mockFolder(t, m)
				m.EXPECT().Dashboards().Return(nil, errors.New("some error")).AnyTimes()
			},
			callMockDst: func(m *mocks.Mockapi) {
				mockNewDatasource(t, m)
				mockNewFolder(t, m)
			},
			callMockCustomClient: func(m *mocks.MockcustomAPI) {},
			expectedError:        errors.New("some error"),
		},
		"migrating dashboards": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockDataSource(t, m)
				mockFolder(t, m)
				mockDashboards(t, m)
			},
			callMockDst: func(m *mocks.Mockapi) {
				mockNewDatasource(t, m)
				mockNewFolder(t, m)
			},
			callMockCustomClient: func(m *mocks.MockcustomAPI) {
				m.EXPECT().AllRuleGroups().Return(
					nil,
					nil,
				).AnyTimes()
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// redirecting logs to a string buffer instead of stdout for tests
			b := &strings.Builder{}
			log.DiagnosticWriter = b

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// mocks gapi.Client for src and dst grafana endpoints
			mockSrc := mocks.NewMockapi(ctrl)
			mockDst := mocks.NewMockapi(ctrl)
			// mocks grafana.Client for custom grafana client
			mockcustomAPI := mocks.NewMockcustomAPI(ctrl)

			tc.callMockSrc(mockSrc)
			tc.callMockDst(mockDst)
			tc.callMockCustomClient(mockcustomAPI)

			app := App{
				Src: mockSrc,
				Dst: mockDst,
			}

			err := app.Run(CustomGrafanaClient{Client: mockcustomAPI})
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
