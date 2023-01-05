package v1

import (
	"os"

	"github.com/go-playground/validator/v10"
	configparser "github.com/psyb0t/go-config-parser"
)

const (
	configFileEnvVarName = "CONFIGFILE"
	defaultConfigFile    = "./config.yml"
	defaultListenAddress = "0.0.0.0:80"
	defaultLogLevel      = "debug"
)

type storageType string

const (
	storageTypeBadgerDB storageType = "badgerDB"
)

type storageBadgerDBConfig struct {
	DSN string `yaml:"dsn"`
}

type storageConfig struct {
	Type     storageType           `yaml:"type" validate:"required"`
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
	ListenAddress string            `yaml:"listenAddress" validate:"hostname_port"`
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
// over defaults and values set in the config file
func newConfig() (config, error) {
	configFile := os.Getenv(configFileEnvVarName)
	if configFile == "" {
		configFile = defaultConfigFile
	}

	defaults := map[string]interface{}{
		"listenAddress": defaultListenAddress,
		"logger": map[string]interface{}{
			"level":  "debug",
			"format": "json",
		},
		"telegramBot": map[string]interface{}{
			"token":           "",
			"superuserChatID": 0,
		},
	}

	cfg := config{}
	if err := configparser.Parse(configparser.ConfigFileTypeYAML,
		configFile, &cfg, defaults); err != nil {
		return config{}, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}
