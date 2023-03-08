package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v49/github"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/github-token-renewer/pkg/config"
)

type App struct {
	cfg   config.Config
	gh    *GitHub
	idMap map[string]int64
}

func NewApp(cfg config.Config, gh *GitHub) (*App, error) {
	if gh != nil {
		return &App{gh: gh}, nil
	}
	gh, err := New(cfg)
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, gh: gh, idMap: map[string]int64{}}, nil
}

func (a *App) GenerateInstallationToken(ctx context.Context) (map[int64]string, error) {
	tokenMap := map[int64]string{}
	for _, installation := range a.cfg.Installations {
		id, err := a.GetInstallationID(ctx, installation)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if _, ok := tokenMap[id]; ok {
			logrus.Infof("duplicate id %d. skip", id)
			continue
		}
		token, _, err := a.gh.client.Apps.CreateInstallationToken(
			ctx,
			id,
			&github.InstallationTokenOptions{})
		if err != nil {
			logrus.Error(err)
			continue
		}
		tokenMap[id] = token.GetToken()
		logrus.Infof("id %d token registered", id)
	}
	return tokenMap, nil
}

func (a *App) GetInstallationID(ctx context.Context, installation config.Installation) (int64, error) {
	if installation.ID != 0 {
		return installation.ID, nil
	}
	if installation.Organization != "" {
		if id, ok := a.idMap[installation.Organization]; ok {
			return id, nil
		}
		i, _, err := a.gh.client.Apps.FindOrganizationInstallation(
			ctx,
			installation.Organization,
		)
		if err != nil {
			return -1, err
		}
		a.idMap[installation.Organization] = i.GetID()
		return i.GetID(), nil
	}
	if installation.Repository != "" {
		if id, ok := a.idMap[installation.Repository]; ok {
			return id, nil
		}
		s := strings.Split(installation.Repository, "/")
		if len(s) != 2 {
			return -1, fmt.Errorf("invalid repository. repository must be defined as 'owner/repo'")
		}
		i, _, err := a.gh.client.Apps.FindRepositoryInstallation(ctx, s[0], s[1])
		if err != nil {
			return -1, err
		}
		a.idMap[installation.Repository] = i.GetID()
		return i.GetID(), nil
	}
	if installation.Username != "" {
		if id, ok := a.idMap[installation.Username]; ok {
			return id, nil
		}
		i, _, err := a.gh.client.Apps.FindUserInstallation(ctx, installation.Username)
		if err != nil {
			return -1, err
		}
		a.idMap[installation.Username] = i.GetID()
		return i.GetID(), nil

	}
	return -1, fmt.Errorf("invalid input. id, organization, repository or username must be defined on each installations.")
}
