package types

type AuthFile string

func (c *Config) IsAuthFileSet() bool {
	return c.AuthFile != ""
}
