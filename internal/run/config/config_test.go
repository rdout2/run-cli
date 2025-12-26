package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/config"
)

func TestLoadSave(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the home directory to the temporary directory
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// 1. Test loading a non-existent config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error when loading non-existent config, but got: %v", err)
	}
	if cfg.Region != "" {
		t.Errorf("expected empty region for non-existent config, but got: %s", cfg.Region)
	}

	// 2. Test saving a config
	cfg.Region = "us-central1"
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// 3. Test loading an existing config
	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load existing config: %v", err)
	}
	if loadedCfg.Region != "us-central1" {
		t.Errorf("expected region 'us-central1', but got: %s", loadedCfg.Region)
	}

	// 4. Test saving an updated config
	loadedCfg.Region = "europe-west1"
	if err := loadedCfg.Save(); err != nil {
		t.Fatalf("failed to save updated config: %v", err)
	}

	// 5. Test loading the updated config
	updatedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load updated config: %v", err)
	}
	if updatedCfg.Region != "europe-west1" {
		t.Errorf("expected region 'europe-west1', but got: %s", updatedCfg.Region)
	}
}

func TestGetConfigPath(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the home directory to the temporary directory
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	expectedPath := filepath.Join(tmpDir, ".run.yaml")
	actualPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	if actualPath != expectedPath {
		t.Errorf("expected config path '%s', but got: '%s'", expectedPath, actualPath)
	}
}

func TestCorruptedConfig(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the home directory to the temporary directory
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// Create a corrupted config file
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	invalidYAML := []byte("region: us-central1\ninvalid-yaml")
	if err := os.WriteFile(configPath, invalidYAML, 0644); err != nil {
		t.Fatalf("failed to write corrupted config file: %v", err)
	}

	// Test loading the corrupted config
	_, err = config.Load()
	if err == nil {
		t.Fatal("expected an error when loading a corrupted config, but got nil")
	}

	if !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Errorf("expected error message to contain 'cannot unmarshal', but got: %s", err.Error())
	}
}
