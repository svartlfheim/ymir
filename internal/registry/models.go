package registry

import (
	"fmt"
	"strings"
)

type VersionStatus string

type statusesContainer struct {
	Pending   VersionStatus
	Preparing VersionStatus
	Ready     VersionStatus
	Failed    VersionStatus
	Archived  VersionStatus
}

var VersionStatuses statusesContainer = statusesContainer{
	Pending:   "pending",
	Preparing: "preparing",
	Ready:     "ready",
	Failed:    "failed",
	Archived:  "archived",
}

type ModuleVersionFQN struct {
	ModuleFQN `validate:"required"`
	Version   string `validate:"required"`
}

func (mv ModuleVersionFQN) String() string {
	return fmt.Sprintf("%s@%s", mv.ModuleFQN.String(), mv.Version)
}

type ModuleVersion struct {
	Id            string        `json:"id"`
	Version       string        `json:"version"`
	ModuleId      string        `json:"module_id"`
	Source        string        `json:"source"`
	DownloadURL   string        `json:"downloadURL"`
	RepositoryURL string        `json:"repositoryURL"`
	Status        VersionStatus `json:"status"`
}

type ModuleFQN struct {
	Name      string `json:"name" validate:"required"`
	Namespace string `json:"namespace" validate:"required"`
	Provider  string `json:"provider" validate:"required"`
}

func (m ModuleFQN) String() string {
	return fmt.Sprintf("%s/%s/%s", m.Provider, m.Namespace, m.Name)
}

type Module struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Provider  string `json:"provider"`
}

type ModuleFilters struct {
	Provider  string
	Namespace string
}

func BuildModuleTable(mods []Module) (h []string, r [][]string) {
	h = []string{"Provider", "Namespace", "ID", "Name"}

	for _, m := range mods {
		r = append(r, []string{
			m.Provider,
			m.Namespace,
			m.Id,
			m.Name,
		})
	}

	return
}

func BuildModuleVersionsTable(mvs []ModuleVersion) (h []string, r [][]string) {
	h = []string{"Module ID", "ID", "Version", "Source", "Repository URL", "Download URL", "Status"}

	for _, m := range mvs {
		r = append(r, []string{
			m.ModuleId,
			m.Id,
			m.Version,
			m.Source,
			m.RepositoryURL,
			m.DownloadURL,
			string(m.Status),
		})
	}

	return
}

func ParseModuleFQN(s string) (ModuleFQN, error) {
	if s == "" {
		return ModuleFQN{}, ErrCouldNotParseModuleFQN{
			Value:   s,
			Message: "value is empty",
		}
	}

	parts := strings.Split(s, "/")

	if len(parts) != 3 {
		return ModuleFQN{}, ErrCouldNotParseModuleFQN{
			Value:   s,
			Message: "invalid format, expected {provider}/{namespace}/{name}",
		}
	}

	return ModuleFQN{
		Provider:  parts[0],
		Namespace: parts[1],
		Name:      parts[2],
	}, nil
}

func ParseModuleVersionFQN(s string) (ModuleVersionFQN, error) {
	if s == "" {
		return ModuleVersionFQN{}, ErrCouldNotParseModuleVersionFQN{
			Value:   s,
			Message: "value is empty",
		}
	}

	parts := strings.Split(s, "@")

	if len(parts) != 2 {
		return ModuleVersionFQN{}, ErrCouldNotParseModuleVersionFQN{
			Value:   s,
			Message: "invalid format, expected {provider}/{namespace}/{name}",
		}
	}

	modFQN, err := ParseModuleFQN(parts[0])

	if err != nil {
		return ModuleVersionFQN{}, ErrCouldNotParseModuleVersionFQN{
			Value:   s,
			Message: "invalid format, expected {provider}/{namespace}/{name}",
		}
	}

	return ModuleVersionFQN{
		ModuleFQN: modFQN,
		Version:   parts[1],
	}, nil
}
