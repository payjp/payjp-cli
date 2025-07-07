package profiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddProfile(t *testing.T) {
	profiles := &Profiles{
		Profiles: make(map[string]*Profile),
	}

	profile := &Profile{
		Name:              "test",
		TestModeSecretKey: "sk_test_123",
		BaseURL:           "https://api.pay.jp",
		GrpcServerAddress: "cli.pay.jp",
	}

	profiles.AddProfile(profile)

	if len(profiles.Profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(profiles.Profiles))
	}

	if profiles.Profiles["test"] != profile {
		t.Error("Profile was not added correctly")
	}
}

func TestRemoveProfile(t *testing.T) {
	profiles := &Profiles{
		Profiles: map[string]*Profile{
			"test": {
				Name:              "test",
				TestModeSecretKey: "sk_test_123",
			},
		},
	}

	profiles.RemoveProfile("test")

	if len(profiles.Profiles) != 0 {
		t.Errorf("Expected 0 profiles, got %d", len(profiles.Profiles))
	}
}

func TestLoadProfile(t *testing.T) {
	testProfile := &Profile{
		Name:              "test",
		TestModeSecretKey: "sk_test_123",
		BaseURL:           "https://api.pay.jp",
		GrpcServerAddress: "cli.pay.jp",
	}

	profiles := &Profiles{
		Profiles: map[string]*Profile{
			"test": testProfile,
		},
	}

	// Test loading existing profile
	loaded := profiles.LoadProfile("test")
	if loaded == nil {
		t.Fatal("Failed to load existing profile")
	}
	if loaded != testProfile {
		t.Error("Loaded profile does not match original")
	}

	// Test loading non-existing profile
	notFound := profiles.LoadProfile("nonexistent")
	if notFound != nil {
		t.Error("Expected nil for non-existing profile")
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "profiles_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test loading from non-existent file
	configPath := filepath.Join(tempDir, "config.yaml")
	profiles, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load from non-existent file: %v", err)
	}
	if len(profiles.Profiles) != 0 {
		t.Errorf("Expected empty profiles, got %d", len(profiles.Profiles))
	}
	if profiles.persistPath != configPath {
		t.Errorf("Expected persist path %s, got %s", configPath, profiles.persistPath)
	}

	// Test loading from existing file
	yamlContent := `profiles:
  default:
    name: default
    test_mode_secret_key: sk_test_default
    base_url: https://api.pay.jp
    grpc_server_address: cli.pay.jp
  production:
    name: production
    test_mode_secret_key: sk_test_prod
    base_url: https://api.pay.jp
    grpc_server_address: cli.pay.jp`

	err = os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	profiles, err = LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load from existing file: %v", err)
	}
	if len(profiles.Profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(profiles.Profiles))
	}

	// Verify loaded profiles
	defaultProfile := profiles.LoadProfile("default")
	if defaultProfile == nil {
		t.Fatal("Default profile not found")
	}
	if defaultProfile.Name != "default" {
		t.Errorf("Expected name 'default', got %s", defaultProfile.Name)
	}
	if defaultProfile.TestModeSecretKey != "sk_test_default" {
		t.Errorf("Expected key 'sk_test_default', got %s", defaultProfile.TestModeSecretKey)
	}
}

func TestPersist(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "profiles_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.yaml")

	profiles := &Profiles{
		Profiles: map[string]*Profile{
			"test": {
				Name:              "test",
				TestModeSecretKey: "sk_test_123",
				BaseURL:           "https://api.pay.jp",
				GrpcServerAddress: "cli.pay.jp",
			},
		},
		persistPath: configPath,
	}

	// Persist profiles
	err = profiles.Persist()
	if err != nil {
		t.Fatalf("Failed to persist profiles: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load and verify
	loaded, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load persisted file: %v", err)
	}
	if len(loaded.Profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(loaded.Profiles))
	}

	testProfile := loaded.LoadProfile("test")
	if testProfile == nil {
		t.Fatal("Test profile not found after persist")
	}
	if testProfile.TestModeSecretKey != "sk_test_123" {
		t.Errorf("Expected key 'sk_test_123', got %s", testProfile.TestModeSecretKey)
	}
}
