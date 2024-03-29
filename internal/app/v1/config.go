package v1

import (
	"os"

	"github.com/go-playground/validator/v10"
	configparser "github.com/psyb0t/go-config-parser"
)

const (
	configFileEnvVarName = "CONFIGFILE"
	defaultListenAddress = "0.0.0.0:80"
	defaultLogLevel      = "debug"
	defaultLogFormat     = "json"
)

type storageType string

const (
	storageTypeBadgerDB storageType = "badgerDB"
)

type storageBadgerDBConfig struct {
	DSN string `yaml:"dsn"`
}

type storageConfig struct {
	Type     storageType           `validate:"required" yaml:"type"`
	BadgerDB storageBadgerDBConfig `yaml:"badgerDB"`
}

type telegramBotConfig struct {
	Token           string `yaml:"token"`
	SuperuserChatID int64  `yaml:"superuserChatID"`
}

type loggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type config struct {
	ListenAddress string            `validate:"hostname_port" yaml:"listenAddress"`
	Logger        loggerConfig      `yaml:"logger"`
	TelegramBot   telegramBotConfig `yaml:"telegramBot"`
	Storage       storageConfig     `yaml:"storage"`
}

// newConfig reads and parses the configuration file and returns a config
// struct. It returns an error if there was an issue reading or
// parsing the file or if the config struct field validation fails.
//
// Note: configparser uses viper with AutomaticEnv() meaning that
// env vars such as LOGGER_LEVEL or TELEGRAMBOT_TOKEN will be used
// over defaults and values set in the config file.
func newConfig() (config, error) {
	configFile := os.Getenv(configFileEnvVarName)

	defaults := map[string]interface{}{
		"listenAddress": defaultListenAddress,
		"logger": map[string]interface{}{
			"level":  defaultLogLevel,
			"format": defaultLogFormat,
		},
		"storage": map[string]interface{}{
			"type": storageTypeBadgerDB,
			"badgerDB": map[string]interface{}{
				"dsn": "",
			},
		},
		"telegramBot": map[string]interface{}{
			"token":           "",
			"superuserChatID": 0,
		},
	}

	cfg := config{}
	if err := configparser.Parse(configparser.ConfigFileTypeYAML,
		configFile, &cfg, defaults, ""); err != nil {
		return config{}, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}
