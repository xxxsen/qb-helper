package config

import (
	"encoding/json"
	"os"

	"github.com/xxxsen/common/logger"
)

type QBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type Config struct {
	QBConfig        QBConfig         `json:"qb_config"`
	LogConfig       logger.LogConfig `json:"log_config"`
	BlacklistUa     []string         `json:"blacklist_ua"`
	BlacklistRegion []string         `json:"blacklist_region"`
	BlacklistIP     []string         `json:"blacklist_ip"`
	BlacklistPeerID []string         `json:"blacklist_peer_id"`
	Interval        int              `json:"interval"`
}

func Parse(file string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
