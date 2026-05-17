package service

import (
	"gateway-console/internal/model"
	"gateway-console/internal/repository"

	frameworkconsul "github.com/calmlax/aevons-framework/core/consul"
)

type CatalogService struct {
	repo     *repository.StaticRepository
	registry *frameworkconsul.Registry
}

func NewCatalogService(repo *repository.StaticRepository, registry *frameworkconsul.Registry) *CatalogService {
	return &CatalogService{
		repo:     repo,
		registry: registry,
	}
}

func (s *CatalogService) Overview() model.Overview {
	return model.Overview{
		Routes:    s.repo.Routes(),
		Upstreams: s.Upstreams(),
		Consumers: s.repo.Consumers(),
		Plugins:   s.repo.Plugins(),
		Policies:  s.repo.Policies(),
	}
}

func (s *CatalogService) Routes() []model.Route {
	return s.repo.Routes()
}

func (s *CatalogService) Upstreams() []model.Upstream {
	upstreams, _ := s.resolveUpstreams(false)
	return upstreams
}

func (s *CatalogService) Consumers() []model.Consumer {
	return s.repo.Consumers()
}

func (s *CatalogService) Plugins() []model.PluginDefinition {
	return s.repo.Plugins()
}

func (s *CatalogService) Policies() []model.Policy {
	return s.repo.Policies()
}

func (s *CatalogService) SwaggerSources() []model.SwaggerSource {
	return s.repo.SwaggerSources()
}

func (s *CatalogService) ResolvedUpstreams() ([]model.Upstream, error) {
	return s.resolveUpstreams(true)
}

func (s *CatalogService) resolveUpstreams(strict bool) ([]model.Upstream, error) {
	upstreams := s.repo.Upstreams()
	resolved := make([]model.Upstream, 0, len(upstreams))

	for _, upstream := range upstreams {
		next := upstream
		if next.Discovery == "consul" && s.registry != nil {
			instances, err := s.registry.Discover(next.ServiceName)
			if err != nil {
				if strict {
					return nil, err
				}
				resolved = append(resolved, next)
				continue
			}

			if len(instances) == 0 {
				if strict {
					return nil, ErrNoHealthyInstances{ServiceName: next.ServiceName}
				}
				resolved = append(resolved, next)
				continue
			}

			nodes := make(map[string]int, len(instances))
			for _, instance := range instances {
				nodes[instance.Address+":"+itoa(instance.Port)] = 1
			}
			next.Nodes = nodes
		}
		resolved = append(resolved, next)
	}

	return resolved, nil
}

type ErrNoHealthyInstances struct {
	ServiceName string
}

func (e ErrNoHealthyInstances) Error() string {
	return "no healthy consul instances for service " + e.ServiceName
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	buf := make([]byte, 0, 12)
	for v > 0 {
		buf = append([]byte{byte('0' + v%10)}, buf...)
		v /= 10
	}
	return sign + string(buf)
}
