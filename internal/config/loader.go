package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func LoadConfig(path string) (AppConfig, error) {
	path = strings.TrimSpace(path)

	if path == "" {
		return AppConfig{}, ErrConfigPathRequired
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return AppConfig{}, fmt.Errorf("%w: %s", ErrConfigNotFound, path)
		}
		return AppConfig{}, fmt.Errorf("read config file %q, %w", path, err)
	}

	var config AppConfig

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&config); err != nil {
		return AppConfig{}, fmt.Errorf("decode config file %q: %w", path, err)
	}

	config = NormalizeConfig(config)

	if err := ValidateConfig(config); err != nil {
		return AppConfig{}, fmt.Errorf("validate config file %q: %w", path, err)
	}

	return config, nil
}
