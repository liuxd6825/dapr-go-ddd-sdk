package restapp

type TemporalConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Namespace string `yaml:"namespace"`
	TaskQueue string `yaml:"taskQueue"`
}
