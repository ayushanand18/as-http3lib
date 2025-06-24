package config

import (
	"context"
	"os"
	"strings"
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
		if err == nil {
			loadErr = yaml.Unmarshal(file, &configData)
			return
		}
		loadErr = err
	})
}

func GetValue(ctx context.Context, keyString string) interface{} {
	loadConfig()
	if loadErr != nil {
		return nil
	}

	keys := strings.Split(keyString, ".")
	var current interface{} = configData

	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[key]
		if !ok {
			return nil
		}
	}

	return current
}

func GetString(ctx context.Context, keyString, defaultValue string) string {
	val := GetValue(ctx, keyString)
	if str, ok := val.(string); ok {
		return str
	}
	return defaultValue
}

func GetInt(ctx context.Context, keyString string, defaultValue int) int {
	val := GetValue(ctx, keyString)
	if i, ok := val.(int); ok {
		return i
	}
	return defaultValue
}

func GetBool(ctx context.Context, keyString string, defaultValue bool) bool {
	val := GetValue(ctx, keyString)
	if b, ok := val.(bool); ok {
		return b
	}
	return defaultValue
}

func GetBytes(ctx context.Context, keyString string) []byte {
	val := GetValue(ctx, keyString)
	return []byte(val.(string))
}
