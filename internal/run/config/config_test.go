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
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to the temporary directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

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
	cfg.Project = "my-project"
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
	if loadedCfg.Project != "my-project" {
		t.Errorf("expected project 'my-project', but got: %s", loadedCfg.Project)
	}

	// 4. Test saving an updated config
	loadedCfg.Region = "europe-west1"
	loadedCfg.Project = "other-project"
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
	if updatedCfg.Project != "other-project" {
		t.Errorf("expected project 'other-project', but got: %s", updatedCfg.Project)
	}
}

func TestGetConfigPath(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to the temporary directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

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
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to the temporary directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

	// Create a corrupted config file
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	invalidYAML := []byte("region: us-central1\n  invalid-indent:")
	if err := os.WriteFile(configPath, invalidYAML, 0644); err != nil {
		t.Fatalf("failed to write corrupted config file: %v", err)
	}

	// Test loading the corrupted config
	_, err = config.Load()
	if err == nil {
		t.Fatal("expected an error when loading a corrupted config, but got nil")
	}

	if !strings.Contains(err.Error(), "cannot unmarshal config") {
		t.Errorf("expected error message to contain 'cannot unmarshal config', but got: %s", err.Error())
	}
}

func TestSaveCreatesDir(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to a non-existent subdirectory of the temporary directory
	homeDir := filepath.Join(tmpDir, "nonexistent")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

	// Create and save a config
	cfg := &config.Config{Region: "us-east1"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Check that the config file was created
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created")
	}
}
