package main

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalConfigPath := configPath
	originalVendorDir := vendorDir

	defer func() {
		configPath = originalConfigPath
		vendorDir = originalVendorDir
	}()

	configPath = filepath.Join(tmpDir, "config.json")
	vendorDir = tmpDir

	// Create a test config
	testConfig := &Config{
		Repos: []Repo{
			{
				Name:   "effect",
				URL:    "https://github.com/Effect-TS/effect",
				Branch: "main",
				Path:   filepath.Join(tmpDir, "effect"),
			},
		},
	}

	// Save config
	if err := saveConfig(testConfig); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedConfig, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the loaded config matches
	if len(loadedConfig.Repos) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(loadedConfig.Repos))
	}

	if loadedConfig.Repos[0].Name != "effect" {
		t.Errorf("Expected repo name 'effect', got '%s'", loadedConfig.Repos[0].Name)
	}

	if loadedConfig.Repos[0].URL != "https://github.com/Effect-TS/effect" {
		t.Errorf("Expected URL 'https://github.com/Effect-TS/effect', got '%s'", loadedConfig.Repos[0].URL)
	}
}

func TestLoadConfigEmpty(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalConfigPath := configPath
	originalVendorDir := vendorDir

	defer func() {
		configPath = originalConfigPath
		vendorDir = originalVendorDir
	}()

	configPath = filepath.Join(tmpDir, "nonexistent.json")
	vendorDir = tmpDir

	// Load config when file doesn't exist
	loadedConfig, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Should return empty config
	if len(loadedConfig.Repos) != 0 {
		t.Errorf("Expected 0 repos, got %d", len(loadedConfig.Repos))
	}
}

func TestAddRepo(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigPath := configPath
	originalVendorDir := vendorDir

	defer func() {
		configPath = originalConfigPath
		vendorDir = originalVendorDir
	}()

	configPath = filepath.Join(tmpDir, "config.json")
	vendorDir = tmpDir

	// Create initial config
	initialConfig := &Config{Repos: []Repo{}}
	if err := saveConfig(initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Add a repo
	newRepo := Repo{
		Name:   "svelte",
		URL:    "https://github.com/sveltejs/svelte",
		Branch: "main",
		Path:   filepath.Join(tmpDir, "svelte"),
	}

	config, _ := loadConfig()
	config.Repos = append(config.Repos, newRepo)
	if err := saveConfig(config); err != nil {
		t.Fatalf("Failed to save config after adding repo: %v", err)
	}

	// Verify repo was added
	loadedConfig, _ := loadConfig()
	if len(loadedConfig.Repos) != 1 {
		t.Errorf("Expected 1 repo after add, got %d", len(loadedConfig.Repos))
	}

	if loadedConfig.Repos[0].Name != "svelte" {
		t.Errorf("Expected repo name 'svelte', got '%s'", loadedConfig.Repos[0].Name)
	}
}

func TestRemoveRepo(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigPath := configPath
	originalVendorDir := vendorDir

	defer func() {
		configPath = originalConfigPath
		vendorDir = originalVendorDir
	}()

	configPath = filepath.Join(tmpDir, "config.json")
	vendorDir = tmpDir

	// Create config with multiple repos
	testConfig := &Config{
		Repos: []Repo{
			{
				Name:   "effect",
				URL:    "https://github.com/Effect-TS/effect",
				Branch: "main",
				Path:   filepath.Join(tmpDir, "effect"),
			},
			{
				Name:   "svelte",
				URL:    "https://github.com/sveltejs/svelte",
				Branch: "main",
				Path:   filepath.Join(tmpDir, "svelte"),
			},
		},
	}
	if err := saveConfig(testConfig); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Remove one repo
	config, _ := loadConfig()
	for i, repo := range config.Repos {
		if repo.Name == "svelte" {
			config.Repos = append(config.Repos[:i], config.Repos[i+1:]...)
			break
		}
	}
	if err := saveConfig(config); err != nil {
		t.Fatalf("Failed to save config after removing repo: %v", err)
	}

	// Verify repo was removed
	loadedConfig, _ := loadConfig()
	if len(loadedConfig.Repos) != 1 {
		t.Errorf("Expected 1 repo after removal, got %d", len(loadedConfig.Repos))
	}

	if loadedConfig.Repos[0].Name != "effect" {
		t.Errorf("Expected remaining repo to be 'effect', got '%s'", loadedConfig.Repos[0].Name)
	}
}

func TestConfigMarshaling(t *testing.T) {
	config := &Config{
		Repos: []Repo{
			{
				Name:   "test-repo",
				URL:    "https://github.com/test/repo",
				Branch: "main",
				Path:   "/home/user/.vendor/test-repo",
			},
		},
	}

	// Test marshaling to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Test unmarshaling from JSON
	var loadedConfig Config
	if err := json.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loadedConfig.Repos[0].Name != "test-repo" {
		t.Errorf("Expected repo name 'test-repo', got '%s'", loadedConfig.Repos[0].Name)
	}
}
