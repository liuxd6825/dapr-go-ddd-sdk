package env

type Resource struct {
	Namespace string         `yaml:"namespace" json:"namespace"`
	Name      string         `yaml:"name" json:"name"`
	Type      string         `yaml:"type" json:"type"`
	URI       string         `yaml:"uri" json:"uri"`
	Metadata  map[string]any `yaml:"metadata" json:"metadata"`
}

func initResources(env *Env) {
	if env == nil {
		return
	}
	if env.Resources == nil {
		env.Resources = map[string]*Resource{}
		return
	}

	for k, v := range env.Resources {
		v.Name = k
	}
	return
}
