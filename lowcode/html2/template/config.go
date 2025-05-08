package template

import (
	"strings"
)

type Config struct {
	Repos map[string]Repository `yaml:"repos"`
}

type RepositoryType string

const (
	LocalRepositoryType RepositoryType = "local"
	GiteaRepositoryType RepositoryType = "gitea"
)

type Repository struct {
	Name     string
	Type     RepositoryType `yaml:"type"`
	Url      string         `yaml:"url"`
	Repo     string         `yaml:"repo"`
	Branch   string         `yaml:"branch"`
	User     string         `yaml:"user"`
	Password string         `yaml:"password"`
	Path     string         `yaml:"path"`
}

func (c *Config) Init() {
	for k, p := range c.Repos {
		values := strings.Split(k, "@")
		if len(values) != 2 {
			panic(k + " is not a valid template repository name, it should be in the format of name@type")
		}
		p.Name = values[0]
		p.Type = RepositoryType(values[1])
	}
}

func (r *Repository) Validate() {
	if r.Name == "" {
		panic(r.Name + " is not a valid template repository name, it should be in the format of name@type")
	}
	if r.Type == "" {
		panic(r.Type + " is not a valid template repository type, it should be either local or gitea")
	}
	switch r.Type {
	case LocalRepositoryType:
		if r.Path == "" {
			panic(r.Name + " is a local repository, but the path is not set")
		}
	case GiteaRepositoryType:
		if r.Url == "" {
			panic(r.Name + " is a gitea repository, but the url is not set")
		}
		if r.Repo == "" {
			panic(r.Name + " is a gitea repository, but the repo is not set")
		}
		if r.Branch == "" {
			panic(r.Name + " is a gitea repository, but the branch is not set")
		}
		if r.User == "" {
			panic(r.Name + " is a gitea repository, but the user is not set")
		}
		if r.Password == "" {
			panic(r.Name + " is a gitea repository, but the password is not set")
		}
	}
}
