//go:build !integration && !e2e

// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package config

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
	"github.com/skuethe/grafana-oss-team-sync/internal/flags"
	"github.com/spf13/pflag"
)

var (
	ErrCouldNotSetRequiredVariable = errors.New("could not set required environment variable")
	ErrCouldNotLoadYamlIntoConfig  = errors.New("could not load YAML file into config")
	ErrCouldNotLoadCLIIntoConfig   = errors.New("could not load CLI input into config")
)

func TestGetConfigFilePath(t *testing.T) {

	type addTest struct {
		name         string
		inputenv     string
		inputflag    string
		expectedpath string
		expectederr  error
	}

	var tests = []addTest{
		{"config via env var", "config.yaml", "", "config.yaml", nil},
		{"config via flag", "", "config.yaml", "config.yaml", nil},
		{"override config from env via flag", "config.yaml", "config2.yaml", "config2.yaml", nil},
		{"no config", "", "", "config.yaml", ErrNoConfigFileDefined},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			os.Clearenv()
			if err := os.Setenv(configtypes.ConfigVariable, test.inputenv); err != nil {
				t.Fatal("could not set required environment variables", "variable", configtypes.ConfigVariable, "input", test.inputenv, "error", err)
			}

			flags.Config = test.inputflag

			outputpath, outputerr := getConfigFilePath()
			if !errors.Is(outputerr, test.expectederr) {
				t.Errorf("got error: %v, wanted error: %v", outputerr, test.expectederr)
			}
			if outputpath != nil && *outputpath != test.expectedpath {
				t.Errorf("got path: %q, wanted path: %q", *outputpath, test.expectedpath)
			}
		})
	}

}

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
				t.Fatal(fmt.Errorf("%w: %w", ErrCouldNotLoadCLIIntoConfig, err))
			}

			if output := k.String(test.path); output != test.expected {
				t.Errorf("got %q, wanted %q", output, test.expected)
			}
		})
	}
}

func TestLoadOptionalAuthFile(t *testing.T) {

	type addTest struct {
		name                    string
		inputauthfileconfigpath string
		inputauthfilename       string
		expectedenvvar          string
		expectedenvcontent      string
		expectedconfigpath      string
		expectedconfigcontent   string
	}

	var tests = []addTest{
		{
			"no authfile set",
			configtypes.AuthFileParameter,
			"",
			"GOTS_TEST",
			"",
			"test",
			"",
		},
		{
			"authfile set but file does not exist",
			configtypes.AuthFileParameter,
			"doesnotexist.env",
			"GOTS_TEST",
			"",
			"test",
			"",
		},
		{
			"authfile set with but non-existing entry",
			configtypes.AuthFileParameter,
			"../../test/data/unit-tests_authfile.env",
			"GOTS_DOESNOTEXIST",
			"",
			"doesnotexist",
			"",
		},
		{
			"authfile set with content",
			configtypes.AuthFileParameter,
			"../../test/data/unit-tests_authfile.env",
			"GOTS_TEST",
			"valid",
			"test",
			"valid",
		},
		{
			"authfile set with content authfile via string input",
			"authfile",
			"../../test/data/unit-tests_authfile.env",
			"GOTS_TEST",
			"valid",
			"test",
			"valid",
		},
		{
			"wrong authfile reference set with content",
			"authFile",
			"../../test/data/unit-tests_authfile.env",
			"GOTS_TEST",
			"",
			"test",
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Setup koanf instance
			k := koanf.New(".")
			err := k.Load(confmap.Provider(map[string]interface{}{
				test.inputauthfileconfigpath: test.inputauthfilename,
			}, "."), nil)
			if err != nil {
				t.Fatal("could not load test data into koanf instance", "error", err)
			}

			// Clear OS env's
			os.Clearenv()

			if err := loadOptionalAuthFile(k); err != nil {
				t.Log(err)
			}

			// Validate output
			if test.inputauthfilename == "" {
				// we did not specify an authfile -> env var and configpath should not be set
				if _, outputvarset := os.LookupEnv(test.expectedenvvar); outputvarset {
					t.Errorf("no authfile set, but env var still present: %v", os.Getenv(test.expectedenvvar))
				}
				if outputconfigset := k.Exists(test.expectedconfigpath); outputconfigset {
					t.Errorf("no authfile set, but config path still present: %v", k.String(test.expectedconfigpath))
				}
			} else {
				if outputvar := os.Getenv(test.expectedenvvar); outputvar != test.expectedenvcontent {
					t.Errorf("wrong env var - got: %v; wanted: %v", outputvar, test.expectedenvcontent)
				}
				if outputconfig := k.String(test.expectedconfigpath); outputconfig != test.expectedconfigcontent {
					t.Errorf("wrong config - got: %v; wanted: %v", outputconfig, test.expectedconfigcontent)
				}
			}
		})
	}
}

func TestUnmarshalIntoStruct(t *testing.T) {

	type addTest struct {
		name     string
		input    string
		expected configtypes.Config
	}

	var tests = []addTest{
		{
			"minimal config input",
			"../../test/data/unit-tests_config_minimal.yaml",
			configtypes.Config{
				ConfigFile: "",
				LogLevel:   configtypes.LogLevel(0),
				Source:     configtypes.Source("entraid"),
				AuthFile:   "",
				Features: configtypes.Features{
					DisableFolders:       configtypes.FeaturesDisableFoldersDefault,
					DisableUserSync:      configtypes.FeaturesDisableUsersDefault,
					AddLocalAdminToTeams: configtypes.FeaturesAddLocalAdminToTeamsDefault,
				},
				Grafana: configtypes.Grafana{
					AuthType: configtypes.GrafanaAuthTypeDefault,
					Connection: configtypes.GrafanaConnection{
						Scheme:   configtypes.GrafanaConnectionSchemeDefault,
						Host:     configtypes.GrafanaConnectionHostDefault,
						BasePath: configtypes.GrafanaConnectionBasePathDefault,
						Retry:    configtypes.GrafanaConnectionRetryDefault,
					},
				},
				Teams:   configtypes.Teams{""},
				Folders: nil,
			},
		},
		{
			"full config input with fake values",
			"../../test/data/unit-tests_config_full-fake.yaml",
			configtypes.Config{
				ConfigFile: "",
				LogLevel:   configtypes.LogLevel(1),
				Source:     configtypes.Source("someothersource"),
				AuthFile:   "someauthfile.env",
				Features: configtypes.Features{
					DisableFolders:       true,
					DisableUserSync:      true,
					AddLocalAdminToTeams: false,
				},
				Grafana: configtypes.Grafana{
					AuthType: "someotherauthtype",
					Connection: configtypes.GrafanaConnection{
						Scheme:   "someotherscheme",
						Host:     "someotherhost:3001",
						BasePath: "/someotherpath",
						Retry:    42,
					},
				},
				Teams: configtypes.Teams{
					"somegroup-1",
					"somegroup-2",
				},
				Folders: configtypes.Folders{
					"somefolder1": {
						Title:       "somefolder-1",
						Description: "some description 1",
						Permissions: configtypes.FolderPermissions{
							Teams: map[string]configtypes.GrafanaPermission{
								"somegroup-1": 1,
								"somegroup-2": 2,
							},
						},
					},
					"somefolder2": {
						Title:       "SomeFolder-2",
						Description: "Some Description 2",
						Permissions: configtypes.FolderPermissions{
							Teams: map[string]configtypes.GrafanaPermission{
								"somegroup-1": 2,
								"somegroup-2": 4,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			os.Clearenv()
			k := koanf.New(".")

			if err := os.Setenv("GOTS_CONFIG", test.input); err != nil {
				t.Fatal(fmt.Errorf("%w (%v): %w", ErrCouldNotSetRequiredVariable, "GOTS_CONFIG", err))
			}

			// Load YAML config
			if err := loadYAMLFile(k); err != nil {
				t.Fatal(fmt.Errorf("%w: %w", ErrCouldNotLoadYamlIntoConfig, err))
			}

			// Need to also load flags, as this will set default values
			flags.Load()
			if err := loadCLIParameter(k, flags.Instance); err != nil {
				t.Fatal(fmt.Errorf("%w: %w", ErrCouldNotLoadCLIIntoConfig, err))
			}

			if err := unmarshalIntoStruct(k); err != nil {
				t.Fatal(err)
			}

			if output := *Instance; !cmp.Equal(output, test.expected) {
				t.Errorf("difference: %+v", cmp.Diff(output, test.expected))
			}
		})
	}
}

func TestLoad(t *testing.T) {

	type addTest struct {
		name     string
		input    string
		expected configtypes.Config
	}

	var tests = []addTest{
		{
			"minimal config input",
			"../../test/data/unit-tests_config_minimal.yaml",
			configtypes.Config{
				ConfigFile: "../../test/data/unit-tests_config_minimal.yaml",
				LogLevel:   configtypes.LogLevel(0),
				Source:     configtypes.Source("entraid"),
				AuthFile:   "",
				Features: configtypes.Features{
					DisableFolders:       configtypes.FeaturesDisableFoldersDefault,
					DisableUserSync:      configtypes.FeaturesDisableUsersDefault,
					AddLocalAdminToTeams: configtypes.FeaturesAddLocalAdminToTeamsDefault,
				},
				Grafana: configtypes.Grafana{
					AuthType: configtypes.GrafanaAuthTypeDefault,
					Connection: configtypes.GrafanaConnection{
						Scheme:   configtypes.GrafanaConnectionSchemeDefault,
						Host:     configtypes.GrafanaConnectionHostDefault,
						BasePath: configtypes.GrafanaConnectionBasePathDefault,
						Retry:    configtypes.GrafanaConnectionRetryDefault,
					},
				},
				Teams:   configtypes.Teams{""},
				Folders: nil,
			},
		},
		{
			"full config input with non-default values",
			"../../test/data/unit-tests_config_full-nondefault.yaml",
			configtypes.Config{
				ConfigFile: "../../test/data/unit-tests_config_full-nondefault.yaml",
				LogLevel:   configtypes.LogLevel(1),
				Source:     configtypes.Source("entraid"),
				AuthFile:   "../../test/data/unit-tests_authfile.env",
				Features: configtypes.Features{
					DisableFolders:       true,
					DisableUserSync:      true,
					AddLocalAdminToTeams: false,
				},
				Grafana: configtypes.Grafana{
					AuthType: "token",
					Connection: configtypes.GrafanaConnection{
						Scheme:   "https",
						Host:     "myhost:3001",
						BasePath: "/sub/api",
						Retry:    7,
					},
				},
				Teams: configtypes.Teams{
					"group-unit-1",
					"group-unit-2",
				},
				Folders: configtypes.Folders{
					"folderunit1": {
						Title:       "folder-unit-1",
						Description: "folder for unit test 1",
						Permissions: configtypes.FolderPermissions{
							Teams: map[string]configtypes.GrafanaPermission{
								"group-unit-1": 1,
								"group-unit-2": 2,
							},
						},
					},
					"folderunit2": {
						Title:       "Folder-Unit-2",
						Description: "Folder For Unit Test 2",
						Permissions: configtypes.FolderPermissions{
							Teams: map[string]configtypes.GrafanaPermission{
								"group-unit-1": 2,
								"group-unit-2": 4,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Reset to ensure clean results
			os.Clearenv()
			Instance = &configtypes.Config{}

			if err := os.Setenv("GOTS_CONFIG", test.input); err != nil {
				t.Fatal(fmt.Errorf("%w (%v): %w", ErrCouldNotSetRequiredVariable, "GOTS_CONFIG", err))
			}

			// Need to also load flags, as this will set default values
			flags.Load()

			Load()

			if output := *Instance; !cmp.Equal(output, test.expected) {
				t.Errorf("difference: %+v", cmp.Diff(output, test.expected))
			}
		})
	}
}
