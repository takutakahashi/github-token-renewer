package github

import (
	"net/http"
	"os"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v49/github"
	"github.com/takutakahashi/github-token-renewer/pkg/config"
)

type GitHub struct {
	client *github.Client
}

func New(cfg config.Config) (*GitHub, error) {
	privatePem, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	t, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, privatePem)
	if err != nil {
		return nil, err
	}
	t.BaseURL = cfg.URL
	if cfg.URL == "" || strings.Contains(cfg.URL, "https://github.com") {
		client := github.NewClient(&http.Client{Transport: t})
		return &GitHub{client: client}, nil
	} else {
		client, err := github.NewEnterpriseClient(
			cfg.URL,
			cfg.URL,
			&http.Client{
				Transport: t,
			})
		if err != nil {
			return nil, err
		}
		return &GitHub{client: client}, nil
	}
}
