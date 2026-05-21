package app

import "github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

func (a *App) migrateAlertRules(destFolders *[]Folder) (int, error) {
	log.Info()
	log.Info("Migrating alert rules:")

	ruleGroups, err := a.Src.AllAlertRuleGroups()
	if err != nil {
		return 0, err
	}

	migrated := 0
	for folderTitle, groups := range ruleGroups {
		folderUID := searchFolderUID(destFolders, folderTitle)
		for _, group := range groups {
			for _, rule := range group.Rules {
				log.InfoLightf("Alert rule: %s\n", rule.GrafanaAlert.Title)

				provRule := &ProvisioningAlertRule{
					Title:        rule.GrafanaAlert.Title,
					Condition:    rule.GrafanaAlert.Condition,
					Data:         rule.GrafanaAlert.Data,
					FolderUID:    folderUID,
					RuleGroup:    group.Name,
					For:          rule.For,
					NoDataState:  rule.GrafanaAlert.NoDataState,
					ExecErrState: rule.GrafanaAlert.ExecErrState,
					Annotations:  rule.Annotations,
					Labels:       rule.Labels,
				}

				if err := a.Dst.CreateAlertRule(provRule); err != nil {
					log.Errorf("\terror: %s\n", err)
				} else {
					migrated++
				}
			}
		}
	}
	return migrated, nil
}
