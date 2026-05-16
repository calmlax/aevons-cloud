package service

import (
	"context"
	"strings"

	"aevons-cloud/gateway/console/internal/apisixadmin"
	"aevons-cloud/gateway/console/internal/model"
)

type PublishService struct {
	catalog *CatalogService
	client  *apisixadmin.Client
}

func NewPublishService(catalog *CatalogService, client *apisixadmin.Client) *PublishService {
	return &PublishService{
		catalog: catalog,
		client:  client,
	}
}

func (s *PublishService) Plan() model.PublishPlan {
	return model.PublishPlan{
		Routes:    s.catalog.Routes(),
		Upstreams: s.catalog.Upstreams(),
		Consumers: s.catalog.Consumers(),
		Policies:  s.catalog.Policies(),
	}
}

func (s *PublishService) Snapshot() (apisixadmin.PublishSnapshot, error) {
	plan := s.Plan()
	resolvedUpstreams, err := s.catalog.ResolvedUpstreams()
	if err != nil {
		return apisixadmin.PublishSnapshot{}, err
	}
	plan.Upstreams = resolvedUpstreams
	snapshot := apisixadmin.PublishSnapshot{
		Routes:        make([]apisixadmin.RouteResource, 0, len(plan.Routes)),
		Upstreams:     make([]apisixadmin.UpstreamResource, 0, len(plan.Upstreams)),
		Consumers:     make([]apisixadmin.ConsumerResource, 0, len(plan.Consumers)),
		PluginConfigs: make([]apisixadmin.PluginConfigResource, 0, len(plan.Policies)),
		GlobalRules: []apisixadmin.GlobalRuleResource{
			{
				ID: "default-global-rule",
				Plugins: map[string]any{
					"prometheus": map[string]any{},
				},
			},
		},
	}

	for _, route := range plan.Routes {
		plugins := make(map[string]any, len(route.Filters))
		for _, filter := range route.Filters {
			plugins[filter] = map[string]any{}
		}
		snapshot.Routes = append(snapshot.Routes, apisixadmin.RouteResource{
			ID:         route.ID,
			Name:       route.Name,
			URI:        route.URI,
			Methods:    route.Methods,
			Plugins:    plugins,
			UpstreamID: route.UpstreamID,
			Status:     boolToStatus(route.Enabled),
		})
	}

	for _, upstream := range plan.Upstreams {
		snapshot.Upstreams = append(snapshot.Upstreams, apisixadmin.UpstreamResource{
			ID:     upstream.ID,
			Name:   upstream.ServiceName,
			Type:   normalizeLB(upstream.LoadBalance),
			Nodes:  upstream.Nodes,
			Scheme: "http",
			Pass:   "pass",
		})
	}

	for _, consumer := range plan.Consumers {
		plugins := make(map[string]any, len(consumer.Plugins))
		for _, plugin := range consumer.Plugins {
			plugins[plugin] = map[string]any{}
		}
		snapshot.Consumers = append(snapshot.Consumers, apisixadmin.ConsumerResource{
			Username: consumer.ID,
			Plugins:  plugins,
			Labels: map[string]string{
				"display_name": consumer.Name,
				"status":       consumer.Status,
			},
		})
	}

	for _, policy := range plan.Policies {
		snapshot.PluginConfigs = append(snapshot.PluginConfigs, apisixadmin.PluginConfigResource{
			ID:   policy.ID,
			Desc: "Aevons client resource policy",
			Plugins: map[string]any{
				"client-resource-auth": map[string]any{
					"client_id": policy.ClientID,
					"enabled":   policy.Enabled,
					"resources": policy.Resources,
				},
			},
		})
	}

	return snapshot, nil
}

func (s *PublishService) Publish(ctx context.Context) error {
	if s.client == nil {
		return nil
	}

	snapshot, err := s.Snapshot()
	if err != nil {
		return err
	}
	for _, upstream := range snapshot.Upstreams {
		if err := s.client.PutUpstream(ctx, upstream.ID, upstream); err != nil {
			return err
		}
	}
	for _, route := range snapshot.Routes {
		if err := s.client.PutRoute(ctx, route.ID, route); err != nil {
			return err
		}
	}
	for _, consumer := range snapshot.Consumers {
		if err := s.client.PutConsumer(ctx, consumer.Username, consumer); err != nil {
			return err
		}
	}
	for _, pluginConfig := range snapshot.PluginConfigs {
		if err := s.client.PutPluginConfig(ctx, pluginConfig.ID, pluginConfig); err != nil {
			return err
		}
	}
	for _, globalRule := range snapshot.GlobalRules {
		if err := s.client.PutGlobalRule(ctx, globalRule.ID, globalRule); err != nil {
			return err
		}
	}
	return nil
}

func boolToStatus(enabled bool) int {
	if enabled {
		return 1
	}
	return 0
}

func normalizeLB(strategy string) string {
	switch strings.TrimSpace(strings.ToLower(strategy)) {
	case "round_robin", "roundrobin", "":
		return "roundrobin"
	case "random":
		return "chash"
	default:
		return strategy
	}
}
