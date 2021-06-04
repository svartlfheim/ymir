package archive

import (
	"fmt"
	"net/url"

	"github.com/svartlfheim/ymir/internal/config"
)

type ArchiveRepository interface {
	CreateArchive(u string) (location string, err error)
}

var (
	SourceGithub = "github"
	SourceGitlab = "gitlab"

	GithubHostName = "github.com"
	GitlabHostName = "github.com"
)

type factory struct {
	Config config.Ymir
}

type ArchiveFactory interface {
	New(u string) (ArchiveRepository, error)
}

func BuildFactory(c config.Ymir) *factory {
	return &factory{
		Config: c,
	}
}

func (f *factory) New(u string) (ArchiveRepository, error) {
	src, err := extractSourceFromURL(u)

	if err != nil {
		return nil, err
	}

	return buildRepository(src, f.Config)
}

func buildRepository(t string, c config.Ymir) (ArchiveRepository, error) {
	switch t {
	case SourceGithub:
		return &GithubArchive{
			AccessToken: c.Git.Github.AccessToken,
		}, nil
	default:
		return nil, fmt.Errorf("git archive not immplemented for source '%s'", t)

	}
}

// TODO: Add config for extra gitlab domains; it can be privately hosted
func extractSourceFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	host := parsed.Hostname()

	switch host {
	case GithubHostName:
		return SourceGithub, nil

	case GitlabHostName:
		return SourceGitlab, nil

	default:
		return "", fmt.Errorf("git driver not implemented for '%s'", host)
	}
}
