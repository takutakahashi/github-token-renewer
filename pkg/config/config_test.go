package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid config",
			content: `url: https://github.com
privateKeyPath: /path/to/key.pem
appID: 12345
installations:
  - id: 67890
    organization: test-org
    repository: test-repo
    username: test-user
    output:
      kubernetesSecret:
        secretName: test-secret
        secretNamespace: default
        key: token
`,
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			content: `invalid: [yaml content`,
			wantErr: true,
		},
		{
			name: "minimal config",
			content: `url: ""
privateKeyPath: ""
appID: 0
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			// Test Load function
			cfg, err := Load(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cfg == nil {
				t.Error("Load() returned nil config without error")
			}

			// Validate loaded config for valid cases
			if !tt.wantErr && tt.name == "valid config" {
				if cfg.URL != "https://github.com" {
					t.Errorf("URL = %v, want https://github.com", cfg.URL)
				}
				if cfg.AppID != 12345 {
					t.Errorf("AppID = %v, want 12345", cfg.AppID)
				}
			}
		})
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/non/existent/file.yaml")
	if err == nil {
		t.Error("Load() should return error for non-existent file")
	}
}

func TestConfigStructure(t *testing.T) {
	cfg := Config{
		URL:            "https://github.com",
		PrivateKeyPath: "/path/to/key",
		AppID:          123,
		Installations: []Installation{
			{
				ID:           456,
				Organization: "org",
				Repository:   "repo",
				Username:     "user",
				Output: Output{
					KubernetesSecret: &OutputKubernetesSecret{
						SecretName:      "secret",
						SecretNamespace: "namespace",
						Key:             "key",
					},
				},
			},
		},
	}

	if cfg.URL != "https://github.com" {
		t.Errorf("URL = %v, want https://github.com", cfg.URL)
	}
	if len(cfg.Installations) != 1 {
		t.Errorf("Installations length = %v, want 1", len(cfg.Installations))
	}
	if cfg.Installations[0].ID != 456 {
		t.Errorf("Installation ID = %v, want 456", cfg.Installations[0].ID)
	}
}
