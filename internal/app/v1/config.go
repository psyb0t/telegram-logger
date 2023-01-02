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

type config struct {
	ListenAddress string            `yaml:"listenAddress" validate:"hostname_port"`
	LogLevel      string            `yaml:"logLevel"`
	TelegramBot   telegramBotConfig `yaml:"telegramBot"`
	Storage       storageConfig     `yaml:"storage"`
}

func newConfig() (config, error) {
	configFile := os.Getenv(configFileEnvVarName)
	if configFile == "" {
		configFile = defaultConfigFile
	}

	defaults := map[string]interface{}{
		"listenAddress": defaultListenAddress,
		"logLevel":      defaultLogLevel,
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
