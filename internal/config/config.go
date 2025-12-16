package config

import (
	"os"
	"time"

	"github.com/kamogelosekhukhune777/api-gateway/internal/router"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	APIHost         string        `yaml:"api_host"` //"default:0.0.0.0:3000"`
}

type TransportConfig struct {
	DialTimeout           time.Duration `yaml:"dial_timeout"`
	ResponseHeaderTimeout time.Duration `yaml:"response_header_timeout"`
	KeepAlive             time.Duration `yaml:"keep_alive"`
	MaxIdleConnsPerHost   int           `yaml:"max_idle_conns_per_host"`
}

type RouterConfig struct {
	Services map[string]string `yaml:"services"`
	Routes   []router.Route    `yaml:"routes"`
}

type Config struct {
	Build           string
	TransportConfig `yaml:"transport"`
	RouterConfig    `yaml:"router_config"`
	ServerConfig    `yaml:"server_config"`
}

func LoadConfig(path string) (*Config, error) {
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
