package config

import (
	"strings"

	"bitbucket.org/pharmaeasyteam/goframework/config"
	"github.com/spf13/viper"
)

//TokenizerConfig app configuration
type TokenizerConfig struct {
	Server     config.ServerConfig
	LoadModule LoadModule
}

type LoadModule struct {
	KeysetName1  string
	KeysetName2  string
	KeysetName3  string
	KeysetName4  string
	KeysetValue1 string
	KeysetValue2 string
	KeysetValue3 string
	KeysetValue4 string
}

//InitViper viper initialisation
func InitViper(viper *viper.Viper) {

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	config.InitViper(viper, "server")

	viper.SetDefault("LoadModule.KeysetName1", "")
	viper.SetDefault("LoadModule.KeysetName2", "")
	viper.SetDefault("LoadModule.KeysetName3", "")
	viper.SetDefault("LoadModule.KeysetName4", "")
	viper.SetDefault("LoadModule.KeysetValue1", "")
	viper.SetDefault("LoadModule.KeysetValue2", "")
	viper.SetDefault("LoadModule.KeysetValue3", "")
	viper.SetDefault("LoadModule.KeysetValue4", "")
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
