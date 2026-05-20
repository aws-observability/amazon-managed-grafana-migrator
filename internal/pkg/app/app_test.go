package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestApp_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srcMock := newMockAPI(ctrl)
	dstMock := newMockAPI(ctrl)

	// datasources
	srcMock.EXPECT().DataSources().Return([]DataSource{
		{Name: "Prometheus", Type: "prometheus", URL: "http://prom:9090"},
	}, nil)
	dstMock.EXPECT().CreateDataSource(gomock.Any()).Return(nil)

	// folders
	srcMock.EXPECT().Folders().Return([]Folder{
		{ID: 1, UID: "folder-1", Title: "Team A"},
	}, nil)
	dstMock.EXPECT().CreateFolder("Team A", "folder-1").Return(Folder{ID: 10, UID: "folder-1", Title: "Team A"}, nil)
	dstMock.EXPECT().Folders().Return([]Folder{
		{ID: 10, UID: "folder-1", Title: "Team A"},
	}, nil)

	// dashboards
	srcMock.EXPECT().SearchDashboards().Return([]DashboardSearchResult{
		{UID: "dash-1", Title: "My Dashboard", FolderTitle: "Team A"},
	}, nil)
	srcMock.EXPECT().DashboardByUID("dash-1").Return(&Dashboard{
		Dashboard: map[string]interface{}{"title": "My Dashboard", "id": float64(1)},
	}, nil)
	dstMock.EXPECT().SaveDashboard(gomock.Any()).Return(nil)

	// alert rules
	srcMock.EXPECT().AllAlertRuleGroups().Return(RuleGroupsByFolder{}, nil)

	a := App{Src: srcMock, Dst: dstMock}
	err := a.Run()
	require.NoError(t, err)
}
