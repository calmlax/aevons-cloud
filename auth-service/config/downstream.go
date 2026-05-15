package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// DownstreamConfig 定义 system-service 自己维护的下游依赖配置。
type DownstreamConfig struct {
	LogServiceGRPCTarget string `yaml:"log_service_grpc_target"`
}

type downstreamRoot struct {
	Downstream DownstreamConfig `yaml:"downstream"`
}

// LoadDownstream 从当前服务配置中读取 downstream 段。
func LoadDownstream(configDir, env string) (DownstreamConfig, error) {
	cfg := DownstreamConfig{}

	basePath := filepath.Join(configDir, "config.yaml")
	if err := mergeDownstreamFromFile(basePath, &cfg); err != nil {
		return DownstreamConfig{}, err
	}

	if env != "" {
		overlayPath := filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
		if _, err := os.Stat(overlayPath); err == nil {
			if err := mergeDownstreamFromFile(overlayPath, &cfg); err != nil {
				return DownstreamConfig{}, err
			}
		}
	}

	return cfg, nil
}

func mergeDownstreamFromFile(path string, cfg *DownstreamConfig) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}

	var root downstreamRoot
	if err := yaml.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}

	if root.Downstream.LogServiceGRPCTarget != "" {
		cfg.LogServiceGRPCTarget = root.Downstream.LogServiceGRPCTarget
	}

	return nil
}
