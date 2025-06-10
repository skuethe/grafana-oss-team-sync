package types

import (
	"errors"
)

type Source string

const (
	SourceEntraID Source = "entraid"
)

func (s Source) String() string {
	return string(s)
}

func (c *Config) ValdidateSourcePlugin() error {
	switch c.Source {
	case SourceEntraID:
		return nil
	}

	return errors.New("invalid source: " + c.Source.String())
}

func (c *Config) GetSource() Source {
	return c.Source
}
