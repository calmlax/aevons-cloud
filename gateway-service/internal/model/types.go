package model

type GatewayConfig struct {
	TrustedProxies []string `yaml:"trusted_proxies"`
	TimeoutSeconds int      `yaml:"timeout_seconds"`
	MaxBodyBytes   int64    `yaml:"max_body_bytes"`
}

type SwaggerConfig struct {
	Enabled    bool               `yaml:"enabled"`
	UIEnabled  bool               `yaml:"ui_enabled"`
	AllowedIPs []string           `yaml:"allowed_ips"`
	Docs       []SwaggerDocConfig `yaml:"docs"`
}

type SwaggerDocConfig struct {
	Name      string `yaml:"name"`
	ServiceID string `yaml:"service_id"`
	URL       string `yaml:"url"`
}

type ServiceConfig struct {
	ID                string   `yaml:"id"`
	Name              string   `yaml:"name"`
	Prefix            string   `yaml:"prefix"`
	Discovery         string   `yaml:"discovery"`
	LoadBalance       string   `yaml:"load_balance"`
	PassAccessToken   bool     `yaml:"pass_access_token"`
	ExcludeAuthRoutes []string `yaml:"exclude_auth_routes"`
}

type ClientRuleConfig struct {
	ClientID  string   `yaml:"client_id"`
	Enabled   bool     `yaml:"enabled"`
	Resources []string `yaml:"resources"`
}

type RouteRule struct {
	Method    string
	Pattern   string
	IsPrefix  bool
	RawSource string
}

type ServiceRule struct {
	ID               string
	Name             string
	Prefix           string
	MatchPrefix      string
	Discovery        string
	LoadBalance      string
	PassAccessToken  bool
	ExcludeAuthRules []RouteRule
}

type ClientRule struct {
	ClientID    string
	Enabled     bool
	AllowAll    bool
	ExactRules  map[string]struct{}
	PrefixRules []string
}

type ClientIdentity struct {
	ClientID string `json:"client_id"`
	Source   string `json:"source"`
}

type SwaggerSource struct {
	Name      string `json:"name"`
	Service   string `json:"service"`
	TargetURL string `json:"target_url"`
	ProxyURL  string `json:"proxy_url"`
}

type UserIdentity struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	ClientID    string   `json:"client_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type UserContext struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	ClientID    string   `json:"client_id"`
	Permissions []string `json:"permissions"`
	Roles       any      `json:"roles"`
	Depts       any      `json:"depts"`
}

type RequestContext struct {
	RequestID string
	Service   *ServiceRule
	Client    *ClientIdentity
	User      *UserIdentity
}
