package output

import (
	"testing"

	"github.com/takutakahashi/github-token-renewer/pkg/config"
)

func TestKubernetesStruct(t *testing.T) {
	cfg := config.OutputKubernetesSecret{
		SecretName:      "test-secret",
		SecretNamespace: "test-namespace",
		Key:             "token",
	}

	k := &Kubernetes{
		cfg: cfg,
		c:   nil,
	}

	if k.cfg.SecretName != "test-secret" {
		t.Errorf("SecretName = %v, want test-secret", k.cfg.SecretName)
	}
	if k.cfg.SecretNamespace != "test-namespace" {
		t.Errorf("SecretNamespace = %v, want test-namespace", k.cfg.SecretNamespace)
	}
	if k.cfg.Key != "token" {
		t.Errorf("Key = %v, want token", k.cfg.Key)
	}
}

// Note: NewKubernetes and Output require a valid Kubernetes cluster connection
// and are better suited for integration tests. Here we just validate the struct.
func TestOutputKubernetesSecretConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       config.OutputKubernetesSecret
		wantName  string
		wantNs    string
		wantKey   string
	}{
		{
			name: "complete config",
			cfg: config.OutputKubernetesSecret{
				SecretName:      "my-secret",
				SecretNamespace: "my-namespace",
				Key:             "my-key",
			},
			wantName: "my-secret",
			wantNs:   "my-namespace",
			wantKey:  "my-key",
		},
		{
			name: "empty config",
			cfg: config.OutputKubernetesSecret{
				SecretName:      "",
				SecretNamespace: "",
				Key:             "",
			},
			wantName: "",
			wantNs:   "",
			wantKey:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.SecretName != tt.wantName {
				t.Errorf("SecretName = %v, want %v", tt.cfg.SecretName, tt.wantName)
			}
			if tt.cfg.SecretNamespace != tt.wantNs {
				t.Errorf("SecretNamespace = %v, want %v", tt.cfg.SecretNamespace, tt.wantNs)
			}
			if tt.cfg.Key != tt.wantKey {
				t.Errorf("Key = %v, want %v", tt.cfg.Key, tt.wantKey)
			}
		})
	}
}
