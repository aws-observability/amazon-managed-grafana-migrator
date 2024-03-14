package app

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"
)

// migrateDataSources recreates data sources from source workspace
// and returns the number of data sources created
func (a *App) migrateDataSources() (int, error) {
	log.Info("Migrating data sources:")
	dsx, err := a.Src.DataSources()
	if err != nil {
		return 0, err
	}

	migratedDs := 0

	for _, ds := range dsx {
		log.InfoLightf("Data source: %s\n", ds.Name)
		if _, err := a.Dst.NewDataSource(ds); err != nil {
			log.InfoLightf("\terror: %s\n", err)
			continue
		}
		migratedDs++
	}
	return migratedDs, nil
}
