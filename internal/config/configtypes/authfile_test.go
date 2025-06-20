package configtypes

import "testing"

func TestIsAuthFileSet(t *testing.T) {

	type addTest struct {
		name     string
		input    Config
		expected bool
	}

	var tests = []addTest{
		{"authfile is set", Config{AuthFile: "valid"}, true},
		{"authfile is not set", Config{}, false},
		{"authfile set but empty", Config{AuthFile: ""}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := test.input.IsAuthFileSet(); output != test.expected {
				t.Errorf("got %t, wanted %t", output, test.expected)
			}
		})
	}
}
