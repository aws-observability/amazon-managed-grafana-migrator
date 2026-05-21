package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMigrateFolders(t *testing.T) {
	tests := map[string]struct {
		setupMocks func(src, dst *mockAPI)
		expected   int
		expectErr  bool
	}{
		"migrates folders": {
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().Folders().Return([]Folder{
					{ID: 1, UID: "f-1", Title: "Team A"},
					{ID: 2, UID: "f-2", Title: "Team B"},
				}, nil)
				dst.EXPECT().CreateFolder("Team A", "f-1").Return(Folder{ID: 10, UID: "f-1", Title: "Team A"}, nil)
				dst.EXPECT().CreateFolder("Team B", "f-2").Return(Folder{ID: 11, UID: "f-2", Title: "Team B"}, nil)
				dst.EXPECT().Folders().Return([]Folder{
					{ID: 10, UID: "f-1", Title: "Team A"},
					{ID: 11, UID: "f-2", Title: "Team B"},
				}, nil)
			},
			expected: 2,
		},
		"error listing source folders": {
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().Folders().Return(nil, errors.New("error"))
			},
			expected:  0,
			expectErr: true,
		},
		"partial failure creating folders": {
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().Folders().Return([]Folder{
					{ID: 1, UID: "f-1", Title: "Team A"},
					{ID: 2, UID: "f-2", Title: "Team B"},
				}, nil)
				dst.EXPECT().CreateFolder("Team A", "f-1").Return(Folder{}, errors.New("already exists"))
				dst.EXPECT().CreateFolder("Team B", "f-2").Return(Folder{ID: 11, UID: "f-2", Title: "Team B"}, nil)
				dst.EXPECT().Folders().Return([]Folder{
					{ID: 11, UID: "f-2", Title: "Team B"},
				}, nil)
			},
			expected: 1,
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
			resp, err := a.migrateFolders()

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, len(resp.MigratedFolders))
			}
		})
	}
}
