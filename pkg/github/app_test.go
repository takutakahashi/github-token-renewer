package github

import (
	"context"
	"testing"

	"github.com/takutakahashi/github-token-renewer/pkg/config"
)

func TestNewApp(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		gh      *GitHub
		wantErr bool
	}{
		{
			name: "with existing GitHub client",
			cfg:  config.Config{},
			gh: &GitHub{
				client: nil,
			},
			wantErr: false,
		},
		{
			name: "without GitHub client - invalid config",
			cfg: config.Config{
				URL:            "https://github.com",
				PrivateKeyPath: "/non/existent/key.pem",
				AppID:          12345,
			},
			gh:      nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApp(tt.cfg, tt.gh)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if app == nil {
					t.Error("NewApp() returned nil app without error")
				}
				if app.gh == nil {
					t.Error("NewApp() returned app with nil GitHub client")
				}
			}
		})
	}
}

func TestGetInstallationID(t *testing.T) {
	app := &App{
		cfg:   config.Config{},
		gh:    &GitHub{client: nil},
		idMap: map[string]int64{},
	}

	tests := []struct {
		name         string
		installation config.Installation
		wantErr      bool
		wantID       int64
		setupCache   bool
	}{
		{
			name: "with ID set",
			installation: config.Installation{
				ID: 12345,
			},
			wantErr: false,
			wantID:  12345,
		},
		{
			name: "organization cached",
			installation: config.Installation{
				Organization: "test-org",
			},
			wantErr:    false, // Should succeed with cached value
			wantID:     67890,
			setupCache: true,
		},
		{
			name: "invalid repository format",
			installation: config.Installation{
				Repository: "invalid-format",
			},
			wantErr: true,
		},
		{
			name: "repository cached",
			installation: config.Installation{
				Repository: "owner/repo",
			},
			wantErr:    false,
			wantID:     99999,
			setupCache: true,
		},
		{
			name:         "no identifier provided",
			installation: config.Installation{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cache for each test
			app.idMap = map[string]int64{}

			if tt.setupCache {
				if tt.installation.Organization != "" {
					app.idMap[tt.installation.Organization] = tt.wantID
				}
				if tt.installation.Repository != "" {
					app.idMap[tt.installation.Repository] = tt.wantID
				}
				if tt.installation.Username != "" {
					app.idMap[tt.installation.Username] = tt.wantID
				}
			}

			ctx := context.Background()
			id, err := app.GetInstallationID(ctx, tt.installation)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstallationID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && id != tt.wantID {
				t.Errorf("GetInstallationID() = %v, want %v", id, tt.wantID)
			}
		})
	}
}

func TestAppStruct(t *testing.T) {
	app := &App{
		cfg: config.Config{
			AppID: 123,
		},
		gh:    &GitHub{client: nil},
		idMap: map[string]int64{"test": 456},
	}

	if app.cfg.AppID != 123 {
		t.Errorf("AppID = %v, want 123", app.cfg.AppID)
	}

	if len(app.idMap) != 1 {
		t.Errorf("idMap length = %v, want 1", len(app.idMap))
	}

	if app.idMap["test"] != 456 {
		t.Errorf("idMap[test] = %v, want 456", app.idMap["test"])
	}
}
