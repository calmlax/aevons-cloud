package service

import (
	"aevons-cloud/gateway/console/internal/model"
	"aevons-cloud/gateway/console/repository"
)

type CatalogService struct {
	repo *repository.StaticRepository
}

func NewCatalogService(repo *repository.StaticRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) Overview() model.Overview {
	return model.Overview{
		Routes:    s.repo.Routes(),
		Upstreams: s.repo.Upstreams(),
		Consumers: s.repo.Consumers(),
		Plugins:   s.repo.Plugins(),
		Policies:  s.repo.Policies(),
	}
}

func (s *CatalogService) Routes() []model.Route {
	return s.repo.Routes()
}

func (s *CatalogService) Upstreams() []model.Upstream {
	return s.repo.Upstreams()
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
