package app

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

func (a *App) migrateDashboards(destFolders *[]gapi.Folder) (int, error) {
	log.Info()
	log.Info("Migrating dashboards:")
	searchDx, err := a.Src.Dashboards()
	log.Debugf(a.Verbose, "Found %d dashboards in src\n", len(searchDx))
	if err != nil {
		return 0, err
	}

	migratedDashboards := 0

	for _, searchD := range searchDx {
		log.InfoLightf("Dashboard: %s\n", searchD.URL)

		if d, err := a.Src.DashboardByUID(searchD.UID); err == nil {
			folderID := searchFolderID(destFolders, searchD.FolderTitle)
			log.Debugf(a.Verbose,
				"searching Folder ID [src folder ID/UID/Title] [%d/%s/%s]  for dashboard in dst grafana: %d\n",
				searchD.FolderID, searchD.UID, searchD.FolderTitle, folderID)

			newModel := d.Model
			newModel["id"] = nil

			newDashboard := gapi.Dashboard{
				Meta: gapi.DashboardMeta{
					IsStarred: d.Meta.IsStarred,
					Slug:      d.Meta.Slug,
					Folder:    folderID,
				},
				Model:     newModel,
				FolderID:  folderID,
				Overwrite: d.Overwrite,
			}

			if _, err := a.Dst.NewDashboard(newDashboard); err != nil {
				log.Errorf("\terror: %s\n", err)
			} else {
				migratedDashboards++
			}
		} else {
			log.Errorf("\terror: %s\n", err)
		}
	}
	return migratedDashboards, nil
}
