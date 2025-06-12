package configtypes

import (
	"errors"
)

type Source string

const (
	SourceDefault   string = "entraid"
	SourceFlagHelp  string = "the source `plugin` you want to use\nAllowed: entraid"
	SourceFlagShort string = "s"
	SourceParameter string = "source"
	SourceVariable  string = "GOTS_SOURCE"

	SourcePluginEntraID Source = "entraid"
)

func (s Source) String() string {
	return string(s)
}

func (c *Config) ValdidateSourcePlugin() error {
	switch c.Source {
	case SourcePluginEntraID:
		return nil
	}

	return errors.New("invalid source plugin defined: " + c.Source.String())
}

func (c *Config) GetSource() Source {
	return c.Source
}
