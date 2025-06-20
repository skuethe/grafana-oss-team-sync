package configtypes

import (
	"testing"
)

func TestString(t *testing.T) {

	type addTest struct {
		name     string
		input    Source
		expected string
	}

	var tests = []addTest{
		{"existing source", SourcePluginEntraID, "entraid"},
		{"non-existing source", Source("valid"), "valid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.String(); output != test.expected {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}

func TestGetSource(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected Source
	}

	var tests = []addTest{
		{"existing source via const", Config{Source: SourcePluginEntraID}, SourcePluginEntraID},
		{"existing source via string", Config{Source: SourcePluginEntraID}, Source("entraid")},
		{"non-existing source", Config{Source: Source("invalid")}, Source("invalid")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.GetSource(); output != test.expected {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}

func TestValidateSourcePlugin(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected error
	}

	var tests = []addTest{
		{"existing source", Config{Source: SourcePluginEntraID}, nil},
		// {"non-existing source", Config{Source: Source("invalid")}, errors.New("invalid source plugin defined: invalid")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.ValdidateSourcePlugin(); output != test.expected {
				t.Errorf("got %v, wanted %v", output, test.expected)
			}
		})
	}
}
