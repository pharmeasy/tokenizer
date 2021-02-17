package config

import (
	"strings"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/config"
	"github.com/spf13/viper"
)

//TokenizerConfig app configuration
type TokenizerConfig struct {
	Server        config.ServerConfig
	LoadEnvModule LoadEnvModule
}

type LoadEnvModule struct {
	BaseURL        string
	Token          string
	Timeout        time.Duration
	MaxConnections int
	KeysetName     KeysetName
	KeysetValue    KeysetValue
}

type KeysetName struct {
	KeysetName1 string
	KeysetName2 string
	KeysetName3 string
	KeysetName4 string
}

type KeysetValue struct {
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

	viper.SetDefault("LoadEnvModule.BaseURL", "")
	viper.SetDefault("LoadEnvModule.Token", "")
	viper.SetDefault("LoadEnvModule.timeout", 5000*time.Millisecond)
	viper.SetDefault("LoadEnvModule.MaxConnections", 300)

	viper.SetDefault("LoadEnvModule.KeysetName.KeysetName1", "")
	viper.SetDefault("LoadEnvModule.KeysetName.KeysetName2", "")
	viper.SetDefault("LoadEnvModule.KeysetName.KeysetName3", "")
	viper.SetDefault("LoadEnvModule.KeysetName.KeysetName4", "")

	viper.SetDefault("LoadEnvModule.KeysetValue.KeysetValue1", "")
	viper.SetDefault("LoadEnvModule.KeysetValue.KeysetValue2", "")
	viper.SetDefault("LoadEnvModule.KeysetValue.KeysetValue3", "")
	viper.SetDefault("LoadEnvModule.KeysetValue.KeysetValue4", "")
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
