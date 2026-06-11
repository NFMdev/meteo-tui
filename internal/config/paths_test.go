package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveConfigPathUsesDefaultPathWhenCustomPathIsEmpty(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/meteo-config-test")

	got, err := ResolveConfigPath("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join("/tmp/meteo-config-test", AppName, ConfigFileName)

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestResolveConfigPathTrimsCustomPath(t *testing.T) {
	got, err := ResolveConfigPath("   /tmp/meteo/config.json   ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Clean("/tmp/meteo/config.json")

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestResolveConfigPathExpandsHomeDirectory(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	got, err := ResolveConfigPath("~/meteo/config.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join(homeDir, "meteo", "config.json")

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestResolveConfigPathHndlesHomeOnly(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	got, err := ResolveConfigPath("~")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != homeDir {
		t.Fatalf("expected %q, got %q", homeDir, got)
	}
}

func TestDefaultConfigPathUsesAppNameAndConfigFileName(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/meteo-config-test")

	got, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join(AppName, ConfigFileName)

	if !strings.HasSuffix(got, expected) {
		t.Fatalf("expected path end with %q, got %q", expected, got)
	}
}
