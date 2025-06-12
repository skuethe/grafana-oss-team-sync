package configtypes

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
	ConfigDefault   string = ""
	ConfigFlagHelp  string = "the file path to your config (required)"
	ConfigFlagShort string = "c"
	ConfigParameter string = "config"
	ConfigVariable  string = "GOTS_CONFIG"
)
