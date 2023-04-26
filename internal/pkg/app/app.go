package app

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

// allows unit testing to provide mock api clients
type api interface {
	// Alerts(params url.Values) ([]gapi.Alert, error)
	DataSources() ([]*gapi.DataSource, error)
	Dashboards() ([]gapi.FolderDashboardSearchResponse, error)
	DashboardsByIDs(ids []int64) ([]gapi.FolderDashboardSearchResponse, error)
	DashboardByUID(uid string) (*gapi.Dashboard, error)
	Folders() ([]gapi.Folder, error)
	NewDataSource(s *gapi.DataSource) (int64, error)
	NewDashboard(dashboard gapi.Dashboard) (*gapi.DashboardSaveResponse, error)
	NewFolder(title string, uid ...string) (gapi.Folder, error)
	NewAlertRule(ar *gapi.AlertRule) (string, error)
}

// App is the main application struct. Contains all the required clients
type App struct {
	// Grafana api clients for source and destination workspaces
	Src, Dst api
}

// Run orchestrates the migration of grafana contents
func (a *App) Run(srcCustomGrafanaClient CustomGrafanaClient) error {
	log.Info()
	migratedDs, err := a.migrateDataSources()
	if err != nil {
		return err
	}
	log.Success("Migrated ", migratedDs, " data sources")

	fx, err := a.migrateFolders()
	if err != nil {
		return err
	}
	log.Success("Migrated ", len(fx), " folders")

	// TODO: if there's an error on migrateFolders, query dashboard IDs in dst
	// This will ensure that the migration is not partially completed
	dashboards, err := a.migrateDashboards(&fx)
	if err != nil {
		return err
	}
	log.Success("Successfully migrated ", dashboards, " dashboards")
	log.Info()

	alertsMigrated, err := a.migrateAlertRules(fx, srcCustomGrafanaClient)
	if err != nil {
		return err
	}
	log.Success("Successfully migrated ", alertsMigrated, " alerts")
	return nil
}
