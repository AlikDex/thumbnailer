package config

import (
    "log"
    "os"
    "path/filepath"

    "app/pkg/filesystem"

    "gopkg.in/yaml.v3"
)

type StorageConfig struct {
    Path string `yaml:"path"`
}

type ServerConfig struct {
    Port string `yaml:"port"`
}

type Config struct {
    Server  ServerConfig  `yaml:"server"`
    Storage StorageConfig `yaml:"storage"`
    TmpDir  string        `yaml:"tmpDir"`
}

func LoadConfig() *Config {
    configFile := "configs/config.yaml"

    if !filesystem.IsExists(configFile) {
        log.Fatalf("Config file not found: %v", configFile)
    }

    data, err := os.ReadFile(configFile)
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }

    var cfg Config
    err = yaml.Unmarshal(data, &cfg)
    if err != nil {
        log.Fatalf("Failed to parse config file: %v", err)
    }

    cfg.Storage.Path = filepath.Clean(cfg.Storage.Path)

    return &cfg
}
