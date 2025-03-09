package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config struct represents the JSON configuration structure
type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`

	DataKeeper struct {
		HeartbeatTimeout int `json:"heartbeat_timeout"`
	} `json:"data_keeper"`
}

// LoadConfig reads and parses the config.json file
func LoadConfig(filename string) *Config {
	filename = filepath.Join("config", filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	return &cfg
}
