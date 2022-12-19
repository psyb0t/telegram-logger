package v1

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name          string
		configFile    string
		expectError   bool
		expectedValue config
	}{
		{
			name:        "valid config",
			configFile:  "./.fixture/valid-config.yml",
			expectError: false,
			expectedValue: config{
				ListenAddress: "0.0.0.0:80",
				LogLevel:      "info",
			},
		},
		{
			name:          "badly formatted config config",
			configFile:    "./.fixture/badly-formatted-config.yml",
			expectError:   true,
			expectedValue: config{},
		},
		{
			name:          "config containing invalid listen address",
			configFile:    "./.fixture/bad-listen-address-config.yml",
			expectError:   true,
			expectedValue: config{},
		},
		{
			name:          "inexistent config",
			configFile:    "./.fixture/inexistent-config.yml",
			expectError:   true,
			expectedValue: config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(configFileEnvVarName, test.configFile)
			defer os.Unsetenv(configFileEnvVarName)

			actual, err := newConfig()
			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")

					return
				}
			}

			if !reflect.DeepEqual(actual, test.expectedValue) {
				t.Errorf("Expected %v but got %v", test.expectedValue, actual)
			}
		})
	}
}
