package model

type Route struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	URI        string            `json:"uri"`
	Methods    []string          `json:"methods"`
	Predicates []string          `json:"predicates"`
	Filters    []string          `json:"filters"`
	UpstreamID string            `json:"upstream_id"`
	Enabled    bool              `json:"enabled"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Upstream struct {
	ID           string         `json:"id"`
	ServiceName  string         `json:"service_name"`
	Discovery    string         `json:"discovery"`
	LoadBalance  string         `json:"load_balance"`
	HealthPolicy string         `json:"health_policy"`
	Nodes        map[string]int `json:"nodes,omitempty"`
}

type Consumer struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	AuthType string   `json:"auth_type"`
	Status   string   `json:"status"`
	Plugins  []string `json:"plugins"`
}

type PluginDefinition struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Phase       string `json:"phase"`
	Description string `json:"description"`
}

type Policy struct {
	ID        string   `json:"id"`
	ClientID  string   `json:"client_id"`
	Enabled   bool     `json:"enabled"`
	Resources []string `json:"resources"`
}

type SwaggerSource struct {
	Name      string `json:"name"`
	Service   string `json:"service"`
	TargetURL string `json:"target_url"`
	ProxyURL  string `json:"proxy_url,omitempty"`
}

type Overview struct {
	Routes    []Route            `json:"routes"`
	Upstreams []Upstream         `json:"upstreams"`
	Consumers []Consumer         `json:"consumers"`
	Plugins   []PluginDefinition `json:"plugins"`
	Policies  []Policy           `json:"policies"`
}

type PublishPlan struct {
	Routes    []Route    `json:"routes"`
	Upstreams []Upstream `json:"upstreams"`
	Consumers []Consumer `json:"consumers"`
	Policies  []Policy   `json:"policies"`
}
