package config

import (
	"errors"
	"strings"
)

type SourceSchema struct {
	Plugin sourcePlugin
}

type sourcePlugin string

const (
	SourceAzure sourcePlugin = "azure"
)

var Source *SourceSchema

func (s sourcePlugin) String() string {
	return string(s)
}

func valdidateSource() error {
	Source = &SourceSchema{
		Plugin: sourcePlugin(strings.ToLower(K.MustString("source"))),
	}

	switch Source.Plugin {
	case SourceAzure:
		return nil
	}

	return errors.New("invalid source: " + Source.Plugin.String())
}

func GetSource() sourcePlugin {
	return Source.Plugin
}
