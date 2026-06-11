package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config directory %w", err)
	}

	return filepath.Join(configDir, AppName, ConfigFileName), nil
}

func ResolveConfigPath(customPath string) (string, error) {
	customPath = strings.TrimSpace(customPath)

	if customPath == "" {
		return DefaultConfigPath()
	}

	expandedPath, err := expandHomePath(customPath)
	if err != nil {
		return "", err
	}

	return filepath.Clean(expandedPath), nil
}

func expandHomePath(path string) (string, error) {
	if path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve user home directory: %w", err)
		}

		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve user home directory: %w", err)
		}
		return filepath.Join(homeDir, strings.TrimPrefix(path, "~/")), nil
	}

	return path, nil
}
