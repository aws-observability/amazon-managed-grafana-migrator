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
	// SrcInput GrafanaInput
}

// const minAlertingMigrationVersion = 9.4

// Run orchestrates the migration of grafana contents
// TODO: allow specifying what to migrate or all (dashboard, ...)
func (a *App) Run(srcCustomGrafanaClient CustomGrafanaClient) error {
	log.Info()
	migratedDs, err := a.migrateDataSources()
	if err != nil {
		return err
	}
	log.Success("Migrated ", migratedDs, " data sources")

	foldersResponse, err := a.migrateFolders()
	if err != nil {
		return err
	}
	log.Success("Migrated ", len(foldersResponse.MigratedFolders), " folders")

	dashboards, err := a.migrateDashboards(&foldersResponse.SrcFolders)
	if err != nil {
		return err
	}
	log.Success("Migrated ", dashboards, " dashboards")
	log.Info()

	//TODO: restrict alert migration to v9.4 dest
	// if strconv.ParseFloat(srcGrafanaVersion, 64) < minAlertingMigrationVersion {
	// 	log.Debug("Skipping alert migration for version", a.Dst.GrafanaVersion)
	// }

	log.Info("Skipping alert rules migration")
	/*
	alertsMigrated, err := a.migrateAlertRules(fx, srcCustomGrafanaClient)
	if err != nil {
		return err
	}
	log.Success("Migrated ", alertsMigrated, " alerts")
	*/
	return nil
}
