package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMigrateDashboards(t *testing.T) {
	tests := map[string]struct {
		setupMocks func(src, dst *mockAPI)
		folders    []Folder
		expected   int
		expectErr  bool
	}{
		"migrates dashboards": {
			folders: []Folder{{ID: 10, UID: "f-1", Title: "Team A"}},
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().SearchDashboards().Return([]DashboardSearchResult{
					{UID: "d-1", Title: "Dashboard 1", FolderTitle: "Team A"},
				}, nil)
				src.EXPECT().DashboardByUID("d-1").Return(&Dashboard{
					Dashboard: map[string]interface{}{"title": "Dashboard 1", "id": float64(1)},
				}, nil)
				dst.EXPECT().SaveDashboard(gomock.Any()).Return(nil)
			},
			expected: 1,
		},
		"error searching dashboards": {
			folders: []Folder{},
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().SearchDashboards().Return(nil, errors.New("error"))
			},
			expected:  0,
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srcMock := newMockAPI(ctrl)
			dstMock := newMockAPI(ctrl)
			tc.setupMocks(srcMock, dstMock)

			a := App{Src: srcMock, Dst: dstMock}
			count, err := a.migrateDashboards(&tc.folders)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, count)
			}
		})
	}
}
