package github

import (
	"context"

	"github.com/google/go-github/v49/github"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/github-token-renewer/pkg/config"
)

type App struct {
	cfg config.Config
	gh  *GitHub
}

func NewApp(cfg config.Config, gh *GitHub) (*App, error) {
	if gh != nil {
		return &App{gh: gh}, nil
	}
	gh, err := New(cfg)
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, gh: gh}, nil
}

func (a *App) GenerateInstallationToken() (map[int64]string, error) {
	ctx := context.Background()
	ret := map[int64]string{}
	for _, installation := range a.cfg.Installations {
		token, _, err := a.gh.client.Apps.CreateInstallationToken(
			ctx,
			installation.ID,
			&github.InstallationTokenOptions{})
		if err != nil {
			logrus.Error(err)
			continue
		}
		ret[installation.ID] = token.GetToken()
	}
	return ret, nil
}
