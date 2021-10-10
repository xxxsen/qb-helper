package config

import (
	"encoding/json"
	"io/ioutil"
)

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type CronItem struct {
	Name   string      `json:"name"`
	Args   interface{} `json:"args"`
	Enable bool        `json:"enable"`
}

type LogConfig struct {
	File    string `json:"file"`
	Level   string `json:"level"`
	Size    int    `json:"size"`
	Rotate  int    `json:"rotate"`
	KeepDay int    `json:"keep_day"`
	Console bool   `json:"console"`
}

type Config struct {
	Auth     *AuthConfig `json:"auth"`
	CronList []*CronItem `json:"cron_config"`
	Log      LogConfig   `json:"log"`
}

func Parse(file string) (*Config, error) {
	cfg := &Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
