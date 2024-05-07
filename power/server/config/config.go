package config

type JsServerConfig struct {
	Service *Service `yaml:"service"`
	Repos   *Repos   `yaml:"repos"`
}

type Service struct {
	Api  *ApiServer  `yaml:"api"`
	Html *HtmlServer `yaml:"html"`
}

type ApiServer struct {
	Enable *bool  `json:"enable"`
	Url    string `yaml:"url"`
	Path   string `yaml:"path"`
}

type HtmlServer struct {
	Url    string `yaml:"url"`
	Enable *bool  `json:"enable"`
}

type Repos struct {
	Gitea map[string]*GiteaRepo `yaml:"gitea"`
	File  map[string]*FileRepo  `yaml:"file"`
}

type GiteaRepo struct {
	Name     string
	Url      string `yaml:"url"`
	Repo     string `yaml:"repo"`
	Branch   string `yaml:"branch"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type FileRepo struct {
	Name string
	Path string `yaml:"path"`
}

func (c *JsServerConfig) Init() {
	for k, p := range c.Repos.Gitea {
		p.Name = k
	}
	for k, p := range c.Repos.File {
		p.Name = k
	}
}
