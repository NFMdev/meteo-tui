package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteConfig(path string, config AppConfig) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return ErrConfigPathRequired
	}

	config = NormalizeConfig(config)

	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	parentDir := filepath.Dir(path)

	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return fmt.Errorf("create config directory %q: %w", parentDir, err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config file %q: %w", path, err)
	}

	return nil
}
