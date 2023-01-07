package v1

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
				ListenAddress: "0.0.0.0:8080",
				Logger: loggerConfig{
					Level:  "debug",
					Format: "json",
				},
				TelegramBot: telegramBotConfig{
					Token:           "abc",
					SuperuserChatID: 123,
				},
				Storage: storageConfig{
					Type: "badgerDB",
					BadgerDB: storageBadgerDBConfig{
						DSN: "/path/to/db/dir",
					},
				},
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
		{
			name:        "no config file",
			configFile:  "",
			expectError: false,
			expectedValue: config{
				ListenAddress: defaultListenAddress,
				Logger: loggerConfig{
					Level:  defaultLogLevel,
					Format: defaultLogFormat,
				},
				Storage: storageConfig{
					Type: storageTypeBadgerDB,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(configFileEnvVarName, test.configFile)
			defer os.Unsetenv(configFileEnvVarName)

			actual, err := newConfig()
			if test.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, actual)
		})
	}
}
