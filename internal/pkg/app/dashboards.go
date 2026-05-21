package app

import "github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

func (a *App) migrateDashboards(destFolders *[]Folder) (int, error) {
	log.Info()
	log.Info("Migrating dashboards:")
	searchDx, err := a.Src.SearchDashboards()
	if err != nil {
		return 0, err
	}
	log.Debugf(a.Verbose, "Found %d dashboards in src\n", len(searchDx))

	migrated := 0
	for _, searchD := range searchDx {
		log.InfoLightf("Dashboard: %s\n", searchD.URL)

		d, err := a.Src.DashboardByUID(searchD.UID)
		if err != nil {
			log.Errorf("\terror: %s\n", err)
			continue
		}

		folderUID := searchFolderUID(destFolders, searchD.FolderTitle)
		model := d.Dashboard
		model["id"] = nil

		req := DashboardSaveRequest{
			Dashboard: model,
			FolderUID: folderUID,
			Overwrite: true,
		}

		if err := a.Dst.SaveDashboard(req); err != nil {
			log.Errorf("\terror: %s\n", err)
		} else {
			migrated++
		}
	}
	return migrated, nil
}
