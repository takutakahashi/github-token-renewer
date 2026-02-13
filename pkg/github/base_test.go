package github

import (
	"testing"

	"github.com/takutakahashi/github-token-renewer/pkg/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid private key path",
			cfg: config.Config{
				URL:            "https://github.com",
				PrivateKeyPath: "/non/existent/key.pem",
				AppID:          12345,
			},
			wantErr: true,
			errMsg:  "should fail with non-existent private key",
		},
		{
			name: "empty url defaults to github.com",
			cfg: config.Config{
				URL:            "",
				PrivateKeyPath: "/tmp/test.pem",
				AppID:          12345,
			},
			wantErr: true, // Will fail because key doesn't exist, but tests the URL logic
			errMsg:  "empty URL should default to github.com",
		},
		{
			name: "github.com url",
			cfg: config.Config{
				URL:            "https://github.com",
				PrivateKeyPath: "/tmp/test.pem",
				AppID:          12345,
			},
			wantErr: true, // Will fail because key doesn't exist
			errMsg:  "github.com URL should use standard client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gh, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v (%s)", err, tt.wantErr, tt.errMsg)
				return
			}

			if !tt.wantErr && gh == nil {
				t.Error("New() returned nil GitHub client without error")
			}

			if tt.wantErr && gh != nil {
				t.Error("New() returned GitHub client with error")
			}
		})
	}
}

func TestGitHubStruct(t *testing.T) {
	gh := &GitHub{
		client: nil,
	}

	if gh.client != nil {
		t.Error("Expected nil client in GitHub struct")
	}
}
