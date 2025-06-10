package types

type Config struct {
	ConfigFile string   `yaml:"config"`
	LogLevel   LogLevel `yaml:"loglevel"`
	Source     Source   `yaml:"source"`
	AuthFile   AuthFile `yaml:"authfile"`
	Features   Features `yaml:"features"`
	Grafana    Grafana  `yaml:"grafana"`
	Teams      Teams    `yaml:"teams"`
	Folders    Folders  `yaml:"folders"`
}

const (
	ConfigParamAuthFile string = "authfile"
)
