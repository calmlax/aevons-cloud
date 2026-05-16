package repository

import "aevons-cloud/gateway/console/internal/model"

type StaticRepository struct{}

func NewStaticRepository() *StaticRepository {
	return &StaticRepository{}
}

func (r *StaticRepository) Routes() []model.Route {
	return []model.Route{
		{
			ID:         "auth-route",
			Name:       "auth-service route",
			URI:        "/api/v1/auth/*",
			Methods:    []string{"GET", "POST", "PUT", "DELETE"},
			Predicates: []string{"Path=/api/v1/auth/**"},
			Filters:    []string{"jwt-enterprise-auth", "client-resource-auth"},
			UpstreamID: "auth-service-upstream",
			Enabled:    true,
		},
		{
			ID:         "system-route",
			Name:       "system-service route",
			URI:        "/api/v1/system/*",
			Methods:    []string{"GET", "POST", "PUT", "DELETE"},
			Predicates: []string{"Path=/api/v1/system/**"},
			Filters:    []string{"jwt-enterprise-auth", "client-resource-auth"},
			UpstreamID: "system-service-upstream",
			Enabled:    true,
		},
		{
			ID:         "log-route",
			Name:       "log-service route",
			URI:        "/api/v1/log/*",
			Methods:    []string{"GET", "POST", "DELETE"},
			Predicates: []string{"Path=/api/v1/log/**"},
			Filters:    []string{"jwt-enterprise-auth", "audit-log"},
			UpstreamID: "log-service-upstream",
			Enabled:    true,
		},
	}
}

func (r *StaticRepository) Upstreams() []model.Upstream {
	return []model.Upstream{
		{
			ID:           "auth-service-upstream",
			ServiceName:  "auth-service",
			Discovery:    "consul",
			LoadBalance:  "round_robin",
			HealthPolicy: "healthy_only",
			StaticNodesBackup: map[string]int{
				"127.0.0.1:10601": 1,
			},
		},
		{
			ID:           "system-service-upstream",
			ServiceName:  "system-service",
			Discovery:    "consul",
			LoadBalance:  "round_robin",
			HealthPolicy: "healthy_only",
			StaticNodesBackup: map[string]int{
				"127.0.0.1:10702": 1,
			},
		},
		{
			ID:           "log-service-upstream",
			ServiceName:  "log-service",
			Discovery:    "consul",
			LoadBalance:  "round_robin",
			HealthPolicy: "healthy_only",
			StaticNodesBackup: map[string]int{
				"127.0.0.1:10803": 1,
			},
		},
	}
}

func (r *StaticRepository) Consumers() []model.Consumer {
	return []model.Consumer{
		{ID: "admin-web", Name: "Admin Web", AuthType: "jwt", Status: "enabled", Plugins: []string{"client-resource-auth"}},
		{ID: "portal-web", Name: "Portal Web", AuthType: "jwt", Status: "enabled", Plugins: []string{"client-resource-auth"}},
	}
}

func (r *StaticRepository) Plugins() []model.PluginDefinition {
	return []model.PluginDefinition{
		{Name: "client-resource-auth", Category: "security", Phase: "access", Description: "客户端资源访问控制与 ALL 语义约束"},
		{Name: "jwt-enterprise-auth", Category: "security", Phase: "access", Description: "企业 JWT 鉴权、用户上下文注入"},
		{Name: "tenant-isolation", Category: "governance", Phase: "access", Description: "多租户上下文识别与隔离"},
		{Name: "audit-log", Category: "observe", Phase: "log", Description: "审计日志与安全事件上报"},
		{Name: "ai-risk-control", Category: "security", Phase: "access", Description: "预留 AI 风控与 WAF 扩展"},
	}
}

func (r *StaticRepository) Policies() []model.Policy {
	return []model.Policy{
		{
			ID:        "policy-admin-web",
			ClientID:  "admin-web",
			Enabled:   true,
			Resources: []string{"ALL"},
		},
		{
			ID:        "policy-portal-web",
			ClientID:  "portal-web",
			Enabled:   true,
			Resources: []string{"POST:/api/v1/auth/login", "POST:/api/v1/auth/refresh", "GET:/api/v1/auth/user"},
		},
	}
}

func (r *StaticRepository) SwaggerSources() []model.SwaggerSource {
	return []model.SwaggerSource{
		{
			Name:      "auth-service",
			Service:   "auth-service",
			TargetURL: "http://127.0.0.1:10701/api/swagger.json",
		},
		{
			Name:      "system-service",
			Service:   "system-service",
			TargetURL: "http://127.0.0.1:10702/api/swagger.json",
		},
		{
			Name:      "log-service",
			Service:   "log-service",
			TargetURL: "http://127.0.0.1:10703/api/swagger.json",
		},
	}
}
