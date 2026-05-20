package app

import "github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

// FoldersResponse holds both folders in the source and destination Grafana
type FoldersResponse struct {
	SrcFolders      []Folder
	MigratedFolders []Folder
	AllDstFolders   []Folder
}

func (a *App) migrateFolders() (*FoldersResponse, error) {
	log.Info()
	log.Info("Migrating folders:")

	fx, err := a.Src.Folders()
	if err != nil {
		return nil, err
	}
	log.Debugf(a.Verbose, "Source Grafana folders found: %d\n", len(fx))

	var newFx []Folder
	for _, f := range fx {
		log.InfoLightf("Folder: %s\n", f.Title)
		newF, err := a.Dst.CreateFolder(f.Title, f.UID)
		if err != nil {
			log.Errorf("\terror: %s [%s]\n", f.Title, err)
		} else {
			newFx = append(newFx, newF)
		}
	}

	allDstFolders, _ := a.Dst.Folders()
	return &FoldersResponse{fx, newFx, allDstFolders}, nil
}

func searchFolderUID(fx *[]Folder, title string) string {
	for _, f := range *fx {
		if f.Title == title {
			return f.UID
		}
	}
	return ""
}
