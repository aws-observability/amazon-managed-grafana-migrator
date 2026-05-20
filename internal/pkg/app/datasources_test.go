package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMigrateDataSources(t *testing.T) {
	tests := map[string]struct {
		setupMocks func(src, dst *mockAPI)
		expected   int
		expectErr  bool
	}{
		"migrates successfully": {
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().DataSources().Return([]DataSource{
					{Name: "Prometheus", Type: "prometheus"},
					{Name: "CloudWatch", Type: "cloudwatch"},
				}, nil)
				dst.EXPECT().CreateDataSource(gomock.Any()).Return(nil).Times(2)
			},
			expected: 2,
		},
		"skips failed datasources": {
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().DataSources().Return([]DataSource{
					{Name: "Prometheus", Type: "prometheus"},
					{Name: "CloudWatch", Type: "cloudwatch"},
				}, nil)
				dst.EXPECT().CreateDataSource(gomock.Any()).Return(errors.New("conflict"))
				dst.EXPECT().CreateDataSource(gomock.Any()).Return(nil)
			},
			expected: 1,
		},
		"error listing datasources": {
			setupMocks: func(src, dst *mockAPI) {
				src.EXPECT().DataSources().Return(nil, errors.New("connection error"))
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
			count, err := a.migrateDataSources()

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, count)
			}
		})
	}
}
