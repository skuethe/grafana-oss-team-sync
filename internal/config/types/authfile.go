package types

type AuthFile string

const (
	AuthFileParameter string = "authfile"
	AuthFileVariable  string = "GOTS_AUTHFILE"
)

func (c *Config) IsAuthFileSet() bool {
	return c.AuthFile != ""
}
