package config

import (
	"context"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	configData map[string]interface{}
	once       sync.Once
	loadErr    error
)

func loadConfig() {
	once.Do(func() {
		file, err := os.ReadFile("config.yaml")
		if err != nil {
			loadErr = err
			return
		}
		err = yaml.Unmarshal(file, &configData)
		if err != nil {
			loadErr = err
		}
	})
}

func GetString(ctx context.Context, keyString, defaultValue string) string {
	loadConfig()
	if loadErr != nil {
		return defaultValue
	}
	if val, ok := configData[keyString]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func GetBool(ctx context.Context, keyString string, defaultValue bool) bool {
	loadConfig()
	if loadErr != nil {
		return defaultValue
	}
	if val, ok := configData[keyString]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}
