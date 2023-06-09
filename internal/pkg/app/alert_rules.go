// Package app provides the grafana migration logic
package app

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/grafana"
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

type customAPI interface {
	AllRuleGroups() (grafana.RuleGroupsByFolder, error)
}

// CustomGrafanaClient is another grafana client in the repo for methods not implemented by gapi
type CustomGrafanaClient struct {
	Client customAPI
}

// migrateAlerts migrate alerts from src to dst
func (a *App) migrateAlertRules(folders []gapi.Folder, customGrafanaClient CustomGrafanaClient) (int, error) {
	log.Info()
	log.Info("Migrating alerting rules:")
	migratedAlertRules := 0

	ruleGroups, err := customGrafanaClient.Client.AllRuleGroups()
	if err != nil {
		return migratedAlertRules, err
	}

	for folder, ruleGroups := range ruleGroups {
		// search for folder uid from name
		uid := searchFolderUID(&folders, folder)
		// log.Debug("Folder = ", folder, ", UID = ", uid)

		for _, rg := range ruleGroups {
			for _, r := range rg.Rules {
				gapiR := convertAlertRule(r, uid)
				log.Debugf("Alerting rule: %s\n", gapiR.Title)

				_, err := a.Dst.NewAlertRule(&gapiR)
				if err == nil {
					migratedAlertRules++
				} else {
					log.Debug("Error creating alerting rule: ", err)
				}
			}
		}
	}

	return migratedAlertRules, nil
}

func convertAlertRule(rule grafana.AlertRule, folderUID string) gapi.AlertRule {
	ax := make([]*gapi.AlertQuery, 0)

	for _, d := range rule.Alert.Data {
		if rule.Alert.Data != nil {
			ax = append(ax, &gapi.AlertQuery{
				DatasourceUID: d.DatasourceUID,
				Model:         d.Model,
				QueryType:     d.QueryType,
				RefID:         d.RefID,
				RelativeTimeRange: gapi.RelativeTimeRange{
					From: d.RelativeTimeRange.From,
					To:   d.RelativeTimeRange.To,
				},
			})
		}
	}

	return gapi.AlertRule{
		Annotations:  rule.Annotations,
		Condition:    rule.Alert.Condition,
		Data:         ax,
		FolderUID:    folderUID,
		RuleGroup:    rule.Alert.RuleGroup,
		Title:        rule.Alert.Title,
		UID:          rule.Alert.UID,
		Updated:      rule.Alert.Updated,
		For:          rule.For,
		ExecErrState: gapi.ExecErrState(rule.Alert.ExecErrState),
		NoDataState:  gapi.NoDataState(rule.Alert.NoDataState),
	}
}
