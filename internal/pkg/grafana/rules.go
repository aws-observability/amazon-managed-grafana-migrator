package grafana

import "time"

// RuleGroup contains alert rules in a group for a folder
type RuleGroup struct {
	Name     string      `json:"name"`
	Interval string      `json:"interval"`
	Rules    []AlertRule `json:"rules"`
}

// AlertRule contains Grafana alert rule
type AlertRule struct {
	Expr string `json:"expr"`
	For  string `json:"for"`
	// Annotations is a map of key/value pairs
	Annotations map[string]string `json:"annotations"`
	// GrafanaAlert is a struct
	Alert Alert `json:"grafana_alert"`
}

// Alert contains a Grafana alert data
type Alert struct {
	ID        int           `json:"id"`
	OrgID     int           `json:"orgId"`
	Title     string        `json:"title"`
	Condition string        `json:"condition"`
	Data      []*AlertQuery `json:"data"`
	Updated   time.Time     `json:"updated"`
	// IntervalSeconds is in seconds
	IntervalSeconds int    `json:"intervalSeconds"`
	Version         int    `json:"version"`
	UID             string `json:"uid"`
	NamespaceUID    string `json:"namespace_uid"`
	NamespaceID     int    `json:"namespace_id"`
	RuleGroup       string `json:"rule_group"`
	NoDataState     string `json:"no_data_state"`
	ExecErrState    string `json:"exec_err_state"`
}

// AlertQuery contains Grafana alert model
type AlertQuery struct {
	RefID             string `json:"refId"`
	QueryType         string `json:"queryType"`
	RelativeTimeRange struct {
		From time.Duration `json:"from"`
		To   time.Duration `json:"to"`
	}
	DatasourceUID string      `json:"datasourceUid"`
	Model         interface{} `json:"model"`
}

// RuleGroupsByFolder contains Grafana rule groups for a Folder
// type RuleGroupsByFolder map[string][]RuleGroupName
type RuleGroupsByFolder map[string][]RuleGroup

// AllRuleGroups fetches all rules group names in folders
func (c *Client) AllRuleGroups() (RuleGroupsByFolder, error) {

	result := make(RuleGroupsByFolder, 0)
	err := c.get("/api/ruler/grafana/api/v1/rules", &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
