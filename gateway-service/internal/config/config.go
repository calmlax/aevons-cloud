package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gateway-service/internal/model"

	"github.com/goccy/go-yaml"
)

type Settings struct {
	Gateway  model.GatewayConfig
	Swagger  model.SwaggerConfig
	Services []model.ServiceConfig
	Clients  []model.ClientRuleConfig
}

type fileConfig struct {
	Gateway  model.GatewayConfig      `yaml:"gateway"`
	Swagger  model.SwaggerConfig      `yaml:"swagger"`
	Services []model.ServiceConfig    `yaml:"services"`
	Clients  []model.ClientRuleConfig `yaml:"clients"`
}

func Load(configDir, env string) (Settings, error) {
	cfg := Settings{
		Gateway: model.GatewayConfig{
			TrustedProxies: []string{"127.0.0.1"},
			TimeoutSeconds: 15,
			MaxBodyBytes:   10 * 1024 * 1024,
		},
		Swagger: model.SwaggerConfig{
			Enabled:    true,
			UIEnabled:  true,
			AllowedIPs: []string{"127.0.0.1", "::1"},
		},
	}

	if err := mergeFile(filepath.Join(configDir, "config.yaml"), &cfg); err != nil {
		return Settings{}, err
	}

	if env != "" {
		overlayPath := filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
		if _, err := os.Stat(overlayPath); err == nil {
			if err := mergeFile(overlayPath, &cfg); err != nil {
				return Settings{}, err
			}
		}
	}

	if err := validate(cfg); err != nil {
		return Settings{}, err
	}

	return cfg, nil
}

func mergeFile(path string, cfg *Settings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read gateway config %s: %w", path, err)
	}

	var next fileConfig
	if err := yaml.Unmarshal(data, &next); err != nil {
		return fmt.Errorf("parse gateway config %s: %w", path, err)
	}

	if len(next.Gateway.TrustedProxies) > 0 {
		cfg.Gateway.TrustedProxies = next.Gateway.TrustedProxies
	}
	if next.Gateway.TimeoutSeconds > 0 {
		cfg.Gateway.TimeoutSeconds = next.Gateway.TimeoutSeconds
	}
	if next.Gateway.MaxBodyBytes > 0 {
		cfg.Gateway.MaxBodyBytes = next.Gateway.MaxBodyBytes
	}

	cfg.Swagger.Enabled = next.Swagger.Enabled
	cfg.Swagger.UIEnabled = next.Swagger.UIEnabled
	if len(next.Swagger.AllowedIPs) > 0 {
		cfg.Swagger.AllowedIPs = next.Swagger.AllowedIPs
	}
	if len(next.Swagger.Docs) > 0 {
		cfg.Swagger.Docs = next.Swagger.Docs
	}

	if len(next.Services) > 0 {
		cfg.Services = next.Services
	}
	if len(next.Clients) > 0 {
		cfg.Clients = next.Clients
	}

	return nil
}

func validate(cfg Settings) error {
	if len(cfg.Services) == 0 {
		return errors.New("gateway config requires at least one service")
	}

	seenServiceIDs := make(map[string]struct{}, len(cfg.Services))
	for _, service := range cfg.Services {
		if strings.TrimSpace(service.ID) == "" || strings.TrimSpace(service.Name) == "" || strings.TrimSpace(service.Prefix) == "" {
			return errors.New("gateway service requires id, name and prefix")
		}
		if _, exists := seenServiceIDs[service.ID]; exists {
			return fmt.Errorf("duplicate service id: %s", service.ID)
		}
		seenServiceIDs[service.ID] = struct{}{}
	}

	seenClients := make(map[string]struct{}, len(cfg.Clients))
	for _, client := range cfg.Clients {
		clientID := strings.TrimSpace(client.ClientID)
		if clientID == "" {
			return errors.New("gateway client rule requires client_id")
		}
		if _, exists := seenClients[clientID]; exists {
			return fmt.Errorf("duplicate client rule: %s", clientID)
		}
		seenClients[clientID] = struct{}{}
	}

	return nil
}
