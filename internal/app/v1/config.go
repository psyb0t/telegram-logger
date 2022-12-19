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

type config struct {
	ListenAddress string `yaml:"listenAddress" validate:"hostname_port"`
	LogLevel      string `yaml:"logLevel"`
}

func newConfig() (config, error) {
	configFile := os.Getenv(configFileEnvVarName)
	if configFile == "" {
		configFile = defaultConfigFile
	}

	defaults := map[string]interface{}{
		"listenAddress": defaultListenAddress,
		"logLevel":      defaultLogLevel,
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
