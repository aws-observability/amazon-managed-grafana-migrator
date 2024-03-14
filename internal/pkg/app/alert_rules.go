// Package app provides the grafana migration logic
package app

import (
	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/grafana"
)

type customAPI interface {
	AllRuleGroups() (grafana.RuleGroupsByFolder, error)
}

// CustomGrafanaClient is another grafana client in the repo for methods not implemented by gapi
type CustomGrafanaClient struct {
	Client customAPI
}
