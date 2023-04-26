package app

import (
	"amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

func (a *App) migrateDashboards(destFolders *[]gapi.Folder) (int, error) {
	log.Info()
	log.Info("Migrating dashboards:")
	searchDx, err := a.Src.Dashboards()
	if err != nil {
		return 0, err
	}

	migratedDashboards := 0

	for _, searchD := range searchDx {
		log.Debugf("Dashboard: %s\n", searchD.URL)

		if d, err := a.Src.DashboardByUID(searchD.UID); err == nil {
			folderID := searchFolderID(destFolders, searchD.FolderTitle)

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
				log.Debugf("\twarning: %s\n", err)
			} else {
				migratedDashboards++
			}
		} else {
			log.Debugf("\twarning: %s\n", err)
		}
	}
	return migratedDashboards, nil
}
