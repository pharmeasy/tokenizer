package config

import (
	"bitbucket.org/pharmaeasyteam/goframework/config"
	"github.com/spf13/viper"
	"strings"
)

//TokenizerConfig app configuration
type TokenizerConfig struct {
	Server config.ServerConfig
}

//InitViper viper initialisation
func InitViper(viper *viper.Viper) {

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	config.InitViper(viper, "server")
}

//Load Load configuration variables from file
func Load(configFile string) (*TokenizerConfig, error) {

	viper := viper.New()
	InitViper(viper)
	cfg := &TokenizerConfig{}
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
