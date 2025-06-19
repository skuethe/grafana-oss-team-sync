package config

import (
	"os"
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
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
		{"set any env var with correct prefix", "GOTS_TEST", "valid", "test", "valid"},
		{"do not load env var with missing prefix", "TEST", "valid", "test", ""},
		{"translate grafana authtype", "GOTS_AUTHTYPE", "valid", "grafana.authtype", "valid"},
		{"translate grafana connection scheme", "GOTS_SCHEME", "valid", "grafana.connection.scheme", "valid"},
		{"translate grafana connection host", "GOTS_HOST", "valid", "grafana.connection.host", "valid"},
		{"translate grafana connection basepath", "GOTS_BASEPATH", "valid", "grafana.connection.basepath", "valid"},
		{"translate grafana connection retry", "GOTS_RETRY", "valid", "grafana.connection.retry", "valid"},
		{"translate feature addlocaladmintoteams", "GOTS_ADDLOCALADMINTOTEAMS", "valid", "features.addLocalAdminToTeams", "valid"},
		{"translate feature disablefolders", "GOTS_DISABLEFOLDERS", "valid", "features.disableFolders", "valid"},
		{"translate feature disableusersync", "GOTS_DISABLEUSERSYNC", "valid", "features.disableUserSync", "valid"},
		// TODO: improve to actually check type here
		{"load one team", "GOTS_TEAMS", "teamA", "teams", "[teamA]"},
		{"load two teams", "GOTS_TEAMS", "teamA,teamB", "teams", "[teamA teamB]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			k := koanf.New(".")
			os.Clearenv()
			if err := os.Setenv(test.variable, test.input); err != nil {
				t.Fatal("could not set required environment variables", "variable", test.variable, "input", test.input, "error", err)
			}

			if err := loadEnvironmentVariables(k); err != nil {
				t.Fatal("could not load environment variables into config. Error:", err)
			}

			if output := k.String(test.path); output != test.expected {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}

func TestLoadCLIParameter(t *testing.T) {

	type addTest struct {
		name     string
		flag     string
		input    string
		path     string
		expected string
	}

	var tests = []addTest{
		{"set a flag", "test", "valid", "test", "valid"},
		{"translate grafana authtype", "authtype", "valid", "grafana.authtype", "valid"},
		{"translate grafana connection scheme", "scheme", "valid", "grafana.connection.scheme", "valid"},
		{"translate grafana connection host", "host", "valid", "grafana.connection.host", "valid"},
		{"translate grafana connection basepath", "basepath", "valid", "grafana.connection.basepath", "valid"},
		{"translate grafana connection retry", "retry", "valid", "grafana.connection.retry", "valid"},
		{"translate feature addlocaladmintoteams", "addlocaladmintoteams", "valid", "features.addLocalAdminToTeams", "valid"},
		{"translate feature disablefolders", "disablefolders", "valid", "features.disableFolders", "valid"},
		{"translate feature disableusersync", "disableusersync", "valid", "features.disableUserSync", "valid"},
		// TODO: improve to actually check type here
		{"load one team", "teams", "teamA", "teams", "[teamA]"},
		{"load two teams", "teams", "teamA,teamB", "teams", "[teamA teamB]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			k := koanf.New(".")

			fs := pflag.NewFlagSet("grafana-oss-team-sync", pflag.ExitOnError)
			fs.String(test.flag, test.input, "")

			if err := loadCLIParameter(k, fs); err != nil {
				t.Fatal("could not load CLI input into config", "error", err)
			}

			if output := k.String(test.path); output != test.expected {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}
