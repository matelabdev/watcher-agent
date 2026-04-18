package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MonitorConfig struct {
	Key      string `yaml:"key"`
	Name     string `yaml:"name"`
	Type     string `yaml:"type"` // tcp | http | rtsp
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	URL      string `yaml:"url"`
	Timeout  int    `yaml:"timeout"`
	Interval int    `yaml:"interval"`
}

type Config struct {
	MasterURL      string          `yaml:"master_url"`
	ProjectToken   string          `yaml:"project_token"`
	ReportInterval int             `yaml:"report_interval"`
	Monitors       []MonitorConfig `yaml:"monitors"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}