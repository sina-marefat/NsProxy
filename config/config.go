package config

import (
	"encoding/json"
	"os"
	"time"
)

var cfg Config

type config struct {
	CacheExpireTTL     string   `json:"cache-expiration-time"`
	ExternalDnsServers []string `json:"external-dns-servers"`
	RedisHost          string   `json:"redis-host"`
	RedisPort          string   `json:"redis-port"`
	ServerHost         string   `json:"server-host"`
	ServerPort         string   `json:"server-port"`
}

type Config struct {
	CacheExpireTTL     time.Duration
	ExternalDnsServers []string
	RedisHost          string
	RedisPort          string
	ServerHost         string
	ServerPort         string
}

func GetConfig() Config {
	return cfg
}

func InitConfig() error {
	var config config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return err
	}
	cfg, err = refineConfig(config)
	return err
}

func refineConfig(cfg config) (Config, error) {
	var refinedCfg Config
	ttl, err := time.ParseDuration(cfg.CacheExpireTTL)
	if err != nil {
		return refinedCfg, err
	}

	refinedCfg = Config{
		CacheExpireTTL:     ttl,
		ExternalDnsServers: cfg.ExternalDnsServers,
		RedisHost:          cfg.RedisHost,
		RedisPort:          cfg.RedisPort,
		ServerHost:         cfg.ServerHost,
		ServerPort:         cfg.ServerPort,
	}

	return refinedCfg, nil
}
