package app

import "github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

func (a *App) migrateDataSources() (int, error) {
	log.Info("Migrating data sources:")
	dsx, err := a.Src.DataSources()
	if err != nil {
		return 0, err
	}

	migrated := 0
	for _, ds := range dsx {
		log.InfoLightf("Data source: %s\n", ds.Name)
		if err := a.Dst.CreateDataSource(&ds); err != nil {
			log.InfoLightf("\terror: %s\n", err)
			continue
		}
		migrated++
	}
	return migrated, nil
}
