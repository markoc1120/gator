package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	tempDir := t.TempDir()
	cases := []struct {
		name           string
		setup          func(path string) error
		expectedConfig *Config
		expectError    bool
	}{
		{
			name: "Missing config file",
			setup: func(path string) error {
				return nil
			},
			expectedConfig: nil,
			expectError:    true,
		},
		{
			name: "Couldn't parse JSON to Go struct",
			setup: func(path string) error {
				return os.WriteFile(filepath.Join(path, ".gatorconfig"), []byte("invalid"), 0644)
			},
			expectedConfig: nil,
			expectError:    true,
		},
		{
			name: "Sucessful reading",
			setup: func(path string) error {
				conf := Config{DBUrl: "postgres://test", CurrentUser: "test_user"}
				data, err := json.Marshal(conf)
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(path, ".gatorconfig.json"), data, 0644)
			},
			expectedConfig: &Config{
				DBUrl:       "postgres://test",
				CurrentUser: "test_user",
			},
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.setup(tempDir)
			if err != nil {
				t.Fatalf("error during setting up the test: %v", err)
			}

			oldHome := os.Getenv("HOME")
			defer os.Setenv("HOME", oldHome)
			os.Setenv("HOME", tempDir)

			conf, err := Read()

			if c.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error but got: %v", err)
				}
				if conf != nil && *conf != *c.expectedConfig {
					t.Errorf("expected config: %v but got: %v", c.expectedConfig, conf)
				}
			}
		})
	}
}

func TestSetUser(t *testing.T) {
	tempDir := t.TempDir()
	initConfig := func(path string) error {
		conf := Config{DBUrl: "postgres://test", CurrentUser: "test_user"}
		data, err := json.Marshal(conf)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(path, ".gatorconfig.json"), data, 0644)
	}

	cases := []struct {
		name           string
		username       string
		setup          func(path string) error
		expectedConfig *Config
		expectError    bool
	}{
		{
			name:     "Sucessful user update",
			username: "new_test_user",
			setup: func(path string) error {
				return initConfig(path)
			},
			expectedConfig: &Config{
				DBUrl:       "postgres://test",
				CurrentUser: "new_test_user",
			},
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.setup(tempDir)
			if err != nil {
				t.Fatalf("error during setting up the test: %v", err)
			}

			oldHome := os.Getenv("HOME")
			defer os.Setenv("HOME", oldHome)
			os.Setenv("HOME", tempDir)

			conf, err := Read()
			if err != nil {
				t.Fatalf("couldn't read the config JSON: %v", err)
			}

			SetUser(c.username, conf)

			conf, err = Read()
			if err != nil {
				t.Fatalf("couldn't read the updated config JSON: %v", err)
			}

			if c.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error but got: %v", err)
				}
				if *c.expectedConfig != *conf {
					t.Errorf("expected config: %v but got: %v", c.expectedConfig, conf)
				}
			}

		})
	}
}
