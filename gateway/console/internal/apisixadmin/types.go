package apisixadmin

type RouteResource struct {
	ID         string                 `json:"id,omitempty"`
	URI        string                 `json:"uri,omitempty"`
	Methods    []string               `json:"methods,omitempty"`
	Hosts      []string               `json:"hosts,omitempty"`
	Vars       []any                  `json:"vars,omitempty"`
	Plugins    map[string]any         `json:"plugins,omitempty"`
	UpstreamID string                 `json:"upstream_id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Status     int                    `json:"status,omitempty"`
	Desc       string                 `json:"desc,omitempty"`
	Labels     map[string]string      `json:"labels,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type UpstreamResource struct {
	ID            string            `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Type          string            `json:"type,omitempty"`
	Nodes         map[string]int    `json:"nodes,omitempty"`
	Scheme        string            `json:"scheme,omitempty"`
	Pass          string            `json:"pass_host,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	ServiceName   string            `json:"service_name,omitempty"`
	DiscoveryType string            `json:"discovery_type,omitempty"`
}

type ConsumerResource struct {
	Username string            `json:"username,omitempty"`
	Plugins  map[string]any    `json:"plugins,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}

type PluginConfigResource struct {
	ID      string            `json:"id,omitempty"`
	Desc    string            `json:"desc,omitempty"`
	Plugins map[string]any    `json:"plugins,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type GlobalRuleResource struct {
	ID      string         `json:"id,omitempty"`
	Plugins map[string]any `json:"plugins,omitempty"`
}

type PublishSnapshot struct {
	Routes        []RouteResource        `json:"routes"`
	Upstreams     []UpstreamResource     `json:"upstreams"`
	Consumers     []ConsumerResource     `json:"consumers"`
	PluginConfigs []PluginConfigResource `json:"plugin_configs"`
	GlobalRules   []GlobalRuleResource   `json:"global_rules"`
}
