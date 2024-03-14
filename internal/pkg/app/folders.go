package app

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

// FoldersResponse holds both folders in the source and destination Grafana
type FoldersResponse struct {
	SrcFolders      []gapi.Folder
	MigratedFolders []gapi.Folder
	AllDstFolders   []gapi.Folder
}

// migrateFolders retrieve folders from source Grafana and use the api to
// create them in the destination. We keep a copy of the source folders
// in case the API fails to create one or more folders (because, it already
// exists for example)
func (a *App) migrateFolders() (*FoldersResponse, error) {
	log.Info()
	log.Info("Migrating folders:")

	fx, err := a.Src.Folders()
	if err != nil {
		return nil, err
	}
	log.Debugf(a.Verbose, "Source Grafana folders found: %d\n", len(fx))
	log.Debug(a.Verbose, fx)

	newFx := []gapi.Folder{}

	for _, f := range fx {
		log.InfoLightf("Folder: %s\n", f.Title)
		newF, err := a.Dst.NewFolder(f.Title, f.UID)
		if err != nil {
			log.Errorf("\terror: %s [%s]\n", f.Title, err)
		} else {
			newFx = append(newFx, newF)
		}
	}

	allDstFolders, _ := a.Dst.Folders()

	return &FoldersResponse{fx, newFx, allDstFolders}, nil
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
