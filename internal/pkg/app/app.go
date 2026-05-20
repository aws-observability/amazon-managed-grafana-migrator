package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws-observability/amazon-managed-grafana-migrator/internal/pkg/log"
)

// Grafana API types

type DataSource struct {
	ID                int64                  `json:"id"`
	UID               string                 `json:"uid"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type"`
	URL               string                 `json:"url"`
	Access            string                 `json:"access"`
	IsDefault         bool                   `json:"isDefault"`
	JSONData          map[string]interface{} `json:"jsonData,omitempty"`
	SecureJSONData    map[string]string      `json:"secureJsonData,omitempty"`
	BasicAuth         bool                   `json:"basicAuth"`
	BasicAuthUser     string                 `json:"basicAuthUser,omitempty"`
	BasicAuthPassword string                 `json:"basicAuthPassword,omitempty"`
}

type Folder struct {
	ID    int64  `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
}

type DashboardSearchResult struct {
	ID          int64  `json:"id"`
	UID         string `json:"uid"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	FolderID    int64  `json:"folderId"`
	FolderUID   string `json:"folderUid"`
	FolderTitle string `json:"folderTitle"`
}

type Dashboard struct {
	Meta      DashboardMeta          `json:"meta"`
	Dashboard map[string]interface{} `json:"dashboard"`
}

type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
	Folder    int64  `json:"folderId"`
	FolderUID string `json:"folderUid"`
}

type DashboardSaveRequest struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	FolderUID string                 `json:"folderUid,omitempty"`
	Overwrite bool                   `json:"overwrite"`
}

type AlertRuleGroup struct {
	Name     string          `json:"name"`
	Interval interface{}     `json:"interval"`
	Rules    []RulerAlertRule `json:"rules,omitempty"`
}

// RulerAlertRule is the format returned by the ruler API (GET)
type RulerAlertRule struct {
	Expr        string            `json:"expr"`
	For         string            `json:"for"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	GrafanaAlert GrafanaAlert     `json:"grafana_alert"`
}

// GrafanaAlert is the nested alert data inside a ruler rule
type GrafanaAlert struct {
	UID          string       `json:"uid"`
	Title        string       `json:"title"`
	Condition    string       `json:"condition"`
	Data         []AlertQuery `json:"data,omitempty"`
	NoDataState  string       `json:"no_data_state"`
	ExecErrState string       `json:"exec_err_state"`
	NamespaceUID string       `json:"namespace_uid"`
	RuleGroup    string       `json:"rule_group"`
}

// ProvisioningAlertRule is the format for the provisioning API (POST)
type ProvisioningAlertRule struct {
	Title        string            `json:"title"`
	Condition    string            `json:"condition"`
	Data         []AlertQuery      `json:"data"`
	FolderUID    string            `json:"folderUID"`
	RuleGroup    string            `json:"ruleGroup"`
	For          string            `json:"for"`
	NoDataState  string            `json:"noDataState"`
	ExecErrState string            `json:"execErrState"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}

type AlertQuery struct {
	RefID             string                 `json:"refId"`
	QueryType         string                 `json:"queryType"`
	RelativeTimeRange map[string]int64       `json:"relativeTimeRange"`
	DatasourceUID     string                 `json:"datasourceUid"`
	Model             map[string]interface{} `json:"model"`
}

type RuleGroupsByFolder map[string][]AlertRuleGroup

// api is the interface for Grafana API operations needed by the migrator
type api interface {
	DataSources() ([]DataSource, error)
	CreateDataSource(ds *DataSource) error
	Folders() ([]Folder, error)
	CreateFolder(title, uid string) (Folder, error)
	SearchDashboards() ([]DashboardSearchResult, error)
	DashboardByUID(uid string) (*Dashboard, error)
	SaveDashboard(req DashboardSaveRequest) error
	AllAlertRuleGroups() (RuleGroupsByFolder, error)
	CreateAlertRule(rule *ProvisioningAlertRule) error
}

// GrafanaClient implements the api interface using HTTP calls
type GrafanaClient struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

// NewGrafanaClient creates a new GrafanaClient
func NewGrafanaClient(baseURL, authToken string, httpClient *http.Client) *GrafanaClient {
	return &GrafanaClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		authToken:  authToken,
		httpClient: httpClient,
	}
}

func (c *GrafanaClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = strings.NewReader(string(b))
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

func (c *GrafanaClient) DataSources() ([]DataSource, error) {
	body, err := c.doRequest(http.MethodGet, "/api/datasources", nil)
	if err != nil {
		return nil, err
	}
	var ds []DataSource
	return ds, json.Unmarshal(body, &ds)
}

func (c *GrafanaClient) CreateDataSource(ds *DataSource) error {
	_, err := c.doRequest(http.MethodPost, "/api/datasources", ds)
	return err
}

func (c *GrafanaClient) Folders() ([]Folder, error) {
	body, err := c.doRequest(http.MethodGet, "/api/folders?limit=1000", nil)
	if err != nil {
		return nil, err
	}
	var folders []Folder
	return folders, json.Unmarshal(body, &folders)
}

func (c *GrafanaClient) CreateFolder(title, uid string) (Folder, error) {
	payload := map[string]string{"title": title, "uid": uid}
	body, err := c.doRequest(http.MethodPost, "/api/folders", payload)
	if err != nil {
		return Folder{}, err
	}
	var f Folder
	return f, json.Unmarshal(body, &f)
}

func (c *GrafanaClient) SearchDashboards() ([]DashboardSearchResult, error) {
	body, err := c.doRequest(http.MethodGet, "/api/search?type=dash-db&limit=5000", nil)
	if err != nil {
		return nil, err
	}
	var results []DashboardSearchResult
	return results, json.Unmarshal(body, &results)
}

func (c *GrafanaClient) DashboardByUID(uid string) (*Dashboard, error) {
	body, err := c.doRequest(http.MethodGet, "/api/dashboards/uid/"+uid, nil)
	if err != nil {
		return nil, err
	}
	var d Dashboard
	return &d, json.Unmarshal(body, &d)
}

func (c *GrafanaClient) SaveDashboard(req DashboardSaveRequest) error {
	_, err := c.doRequest(http.MethodPost, "/api/dashboards/db", req)
	return err
}

func (c *GrafanaClient) AllAlertRuleGroups() (RuleGroupsByFolder, error) {
	body, err := c.doRequest(http.MethodGet, "/api/ruler/grafana/api/v1/rules", nil)
	if err != nil {
		return nil, err
	}
	var result RuleGroupsByFolder
	return result, json.Unmarshal(body, &result)
}

func (c *GrafanaClient) CreateAlertRule(rule *ProvisioningAlertRule) error {
	_, err := c.doRequest(http.MethodPost, "/api/v1/provisioning/alert-rules", rule)
	return err
}

// App is the main application struct
type App struct {
	Src, Dst api
	Verbose  bool
}

// Run orchestrates the migration of grafana contents
func (a *App) Run() error {
	log.Info()
	migratedDs, err := a.migrateDataSources()
	if err != nil {
		return err
	}
	log.Success("Migrated ", migratedDs, " data sources")

	foldersResponse, err := a.migrateFolders()
	log.Debug(a.Verbose, foldersResponse)
	if err != nil {
		return err
	}
	log.Success("Migrated ", len(foldersResponse.MigratedFolders), " folders")

	dashboards, err := a.migrateDashboards(&foldersResponse.AllDstFolders)
	if err != nil {
		return err
	}
	log.Success("Migrated ", dashboards, " dashboards")

	alertRules, err := a.migrateAlertRules(&foldersResponse.AllDstFolders)
	if err != nil {
		log.Errorf("Alert rules migration error: %s\n", err)
	} else {
		log.Success("Migrated ", alertRules, " alert rules")
	}

	log.Info()
	return nil
}
