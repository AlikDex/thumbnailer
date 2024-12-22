package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	Server                 ServerConfig  `yaml:"server"`
	Storage                StorageConfig `yaml:"storage"`
	TmpDir                 string        `yaml:"tmpDir"`
	Upstream               string        `yaml:"upstream"`
	AllovedVideoExtensions []string      `yaml:"allovedVideoExtensions"`
	AllovedImageExtensions []string      `yaml:"allovedImageExtensions"`
	CacheEnable            bool          `yaml:"cacheEnable"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		configFile := "config/config.yaml"

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
		cfg.Upstream = strings.TrimRight(cfg.Upstream, "/")
		cfg.AllovedVideoExtensions = filterEmptyStrings(cfg.AllovedVideoExtensions)

		instance = &cfg
	})

	return instance
}

func filterEmptyStrings(arr []string) []string {
	result := []string{}
	for _, str := range arr {
		if strings.TrimSpace(str) != "" {
			result = append(result, str)
		}
	}

	return result
}
