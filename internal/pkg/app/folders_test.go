package app

import (
	"errors"
	"testing"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/app/mocks"

	"github.com/golang/mock/gomock"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/require"
)

func generateTestFolders(_ *testing.T) []gapi.Folder {
	return []gapi.Folder{
		{
			ID:    1,
			UID:   "uid",
			Title: "test",
			URL:   "http://test.com",
		},
		{
			ID:    2,
			UID:   "uid2",
			Title: "test2",
			URL:   "http://test2.com",
		},
	}
}

func TestApp_migrateFolders(t *testing.T) {
	//test cases
	tests := map[string]struct {
		callMockSrc     func(m *mocks.Mockapi)
		callMockDst     func(m *mocks.Mockapi)
		migratedFolders int
		sourceFolders   int
		expectedError   error
	}{
		"error getting folders from src": {
			callMockSrc: func(m *mocks.Mockapi) {
				m.EXPECT().Folders().Return(nil, errors.New("some error")).AnyTimes()
			},
			callMockDst:     func(m *mocks.Mockapi) {},
			migratedFolders: 0,
			sourceFolders:   0,
			expectedError:   errors.New("some error"),
		},
		"syncing one folder with error": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockFolder(t, m)
			},
			callMockDst: func(m *mocks.Mockapi) {
				mockNewFolderWithError(t, m)
			},
			migratedFolders: 0,
			sourceFolders:   1,
			expectedError:   nil,
		},
		"syncing one folder": {
			callMockSrc: func(m *mocks.Mockapi) {
				mockFolder(t, m)
			},
			callMockDst: func(m *mocks.Mockapi) {
				mockNewFolder(t, m)
			},
			migratedFolders: 1,
			sourceFolders:   1,
			expectedError:   nil,
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

			fx, err := app.migrateFolders()

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Len(t, fx.SrcFolders, tc.sourceFolders)
				require.Len(t, fx.MigratedFolders, tc.migratedFolders)
			}
		})
	}
}

func TestApp_getFolderID(t *testing.T) {

	//test cases
	tests := map[string]struct {
		folders     []gapi.Folder
		folderTitle string
		expectedID  int64
	}{
		"empty list": {
			folders:     []gapi.Folder{},
			folderTitle: "test",
			expectedID:  0,
		},
		"included folder": {
			folders:     generateTestFolders(t),
			folderTitle: "test2",
			expectedID:  2,
		},
		"non included folder": {
			folders:     generateTestFolders(t),
			folderTitle: "test3",
			expectedID:  0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id := searchFolderID(&tc.folders, tc.folderTitle)
			require.Equal(t, id, tc.expectedID)
		})
	}
}

func TestApp_getFolderUID(t *testing.T) {

	//test cases
	tests := map[string]struct {
		folders     []gapi.Folder
		folderTitle string
		expectedUID string
	}{
		"empty list": {
			folders:     []gapi.Folder{},
			folderTitle: "test",
			expectedUID: "",
		},
		"included folder": {
			folders:     generateTestFolders(t),
			folderTitle: "test2",
			expectedUID: "uid2",
		},
		"non included folder": {
			folders: []gapi.Folder{
				{
					ID:    1,
					UID:   "uid",
					Title: "test",
					URL:   "http://test.com",
				},
			},
			folderTitle: "test2",
			expectedUID: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id := searchFolderUID(&tc.folders, tc.folderTitle)
			require.Equal(t, id, tc.expectedUID)
		})
	}
}
