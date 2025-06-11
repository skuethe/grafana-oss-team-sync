package types

type AuthFile string

const (
	AuthFileDefault   string = ""
	AuthFileFlagHelp  string = "an optional file to load authentication data from. File content needs to be in key=value syntax"
	AuthFileFlagShort string = "a"
	AuthFileParameter string = "authfile"
	AuthFileVariable  string = "GOTS_AUTHFILE"
)

func (c *Config) IsAuthFileSet() bool {
	return c.AuthFile != ""
}
