package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Settings struct {
	APISIXAdminURL string `yaml:"apisix_admin_url"`
	APISIXAdminKey string `yaml:"apisix_admin_key"`
	SwaggerUIURL   string `yaml:"swagger_ui_url"`
}

type fileConfig struct {
	Console Settings `yaml:"console"`
}

func Load(configDir, env string) (Settings, error) {
	cfg := Settings{
		APISIXAdminURL: "http://127.0.0.1:9180",
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

	return cfg, nil
}

func mergeFile(path string, cfg *Settings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read console config %s: %w", path, err)
	}

	var next fileConfig
	if err := yaml.Unmarshal(data, &next); err != nil {
		return fmt.Errorf("parse console config %s: %w", path, err)
	}

	if next.Console.APISIXAdminURL != "" {
		cfg.APISIXAdminURL = next.Console.APISIXAdminURL
	}
	if next.Console.APISIXAdminKey != "" {
		cfg.APISIXAdminKey = next.Console.APISIXAdminKey
	}
	if next.Console.SwaggerUIURL != "" {
		cfg.SwaggerUIURL = next.Console.SwaggerUIURL
	}

	return nil
}
