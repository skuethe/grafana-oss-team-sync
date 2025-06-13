package config

import (
	"os"
	"testing"

	"github.com/knadh/koanf/v2"
)

func TestLoadEnvironmentVariables(t *testing.T) {

	type addTest struct {
		name     string
		variable string
		input    string
		path     string
		expected string
	}

	var tests = []addTest{
		{"load env", "GOTS_TEST", "valid", "test", "valid"},
		{"do not load env without prefix", "TEST", "valid", "test", ""},
		{"load one team", "GOTS_TEAMS", "teamA", "teams", "[teamA]"},
		// TODO: improve to actually check type here
		{"load two teams", "GOTS_TEAMS", "teamA,teamB", "teams", "[teamA teamB]"},
		{"translate authtype", "GOTS_AUTHTYPE", "valid", "grafana.authtype", "valid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			k := koanf.New(".")
			os.Clearenv()
			os.Setenv(test.variable, test.input)

			loadEnvironmentVariables(k)

			if output := k.String(test.path); output != test.expected {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}
