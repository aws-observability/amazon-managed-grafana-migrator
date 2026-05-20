package app

import (
	"go.uber.org/mock/gomock"
)

// mockAPI is a mock of the api interface for testing.
type mockAPI struct {
	ctrl     *gomock.Controller
	recorder *mockAPIMockRecorder
}

type mockAPIMockRecorder struct {
	mock *mockAPI
}

func newMockAPI(ctrl *gomock.Controller) *mockAPI {
	mock := &mockAPI{ctrl: ctrl}
	mock.recorder = &mockAPIMockRecorder{mock}
	return mock
}

func (m *mockAPI) EXPECT() *mockAPIMockRecorder {
	return m.recorder
}

func (m *mockAPI) DataSources() ([]DataSource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DataSources")
	ret0, _ := ret[0].([]DataSource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *mockAPIMockRecorder) DataSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "DataSources")
}

func (m *mockAPI) CreateDataSource(ds *DataSource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDataSource", ds)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *mockAPIMockRecorder) CreateDataSource(ds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "CreateDataSource", ds)
}

func (m *mockAPI) Folders() ([]Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Folders")
	ret0, _ := ret[0].([]Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *mockAPIMockRecorder) Folders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "Folders")
}

func (m *mockAPI) CreateFolder(title, uid string) (Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFolder", title, uid)
	ret0, _ := ret[0].(Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *mockAPIMockRecorder) CreateFolder(title, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "CreateFolder", title, uid)
}

func (m *mockAPI) SearchDashboards() ([]DashboardSearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchDashboards")
	ret0, _ := ret[0].([]DashboardSearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *mockAPIMockRecorder) SearchDashboards() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "SearchDashboards")
}

func (m *mockAPI) DashboardByUID(uid string) (*Dashboard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DashboardByUID", uid)
	ret0, _ := ret[0].(*Dashboard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *mockAPIMockRecorder) DashboardByUID(uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "DashboardByUID", uid)
}

func (m *mockAPI) SaveDashboard(req DashboardSaveRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDashboard", req)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *mockAPIMockRecorder) SaveDashboard(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "SaveDashboard", req)
}

func (m *mockAPI) AllAlertRuleGroups() (RuleGroupsByFolder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllAlertRuleGroups")
	ret0, _ := ret[0].(RuleGroupsByFolder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *mockAPIMockRecorder) AllAlertRuleGroups() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "AllAlertRuleGroups")
}

func (m *mockAPI) CreateAlertRule(rule *ProvisioningAlertRule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAlertRule", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *mockAPIMockRecorder) CreateAlertRule(rule any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCall(mr.mock, "CreateAlertRule", rule)
}
