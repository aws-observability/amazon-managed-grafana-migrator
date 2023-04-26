package app

import (
	"amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

func (a *App) migrateFolders() ([]gapi.Folder, error) {
	log.Info()
	log.Info("Migrating folders:")

	fx, err := a.Src.Folders()
	if err != nil {
		return nil, err
	}

	newFx := []gapi.Folder{}

	for _, f := range fx {
		log.Debugf("Folder: %s\n", f.Title)
		newF, err := a.Dst.NewFolder(f.Title, f.UID)
		if err != nil {
			log.Debugf("\terror: %s [%s]\n", f.Title, err)
		}
		newFx = append(newFx, newF)
	}
	return newFx, nil
}

// From a list of folders, get a folder ID (used with destination folders)
func searchFolderID(fx *[]gapi.Folder, title string) int64 {
	for _, f := range *fx {
		if f.Title == title {
			return f.ID
		}
	}
	return 0
}

// From a list of folders, get a folder UID (used with destination folders)
func searchFolderUID(fx *[]gapi.Folder, title string) string {
	for _, f := range *fx {
		if f.Title == title {
			return f.UID
		}
	}
	return ""
}
