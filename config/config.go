package config

import (
	"strings"

	"bitbucket.org/pharmaeasyteam/goframework/config"
	"github.com/spf13/viper"
)

//TokenizerConfig app configuration
type TokenizerConfig struct {
	Server      config.ServerConfig
	VaultModule VaultModule
}

// VaultModule is used to load env variables from vault
type VaultModule struct {
	KeysetConfig KeysetConfig
	KMSConfig    KMSConfig
	DynamoConfig DynamoConfig
	TokenConfig  TokenConfig
}

// KeysetConfig is used to set keysets from vault
type KeysetConfig struct {
	KeysetName1  string
	KeysetName2  string
	KeysetName3  string
	KeysetName4  string
	KeysetValue1 string
	KeysetValue2 string
	KeysetValue3 string
	KeysetValue4 string
}

// KMSConfig is used to set kms arn from vault
type KMSConfig struct {
	AWSKMSKey string
}

// DynamoConfig is used to set dynamodb tablename from vault
type DynamoConfig struct {
	DynamoDBTableName string
}

// TokenConfig is used to set instance id from vault
type TokenConfig struct {
	InstanceID string
}

//InitViper viper initialisation
func InitViper(viper *viper.Viper) {

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	config.InitViper(viper, "server")

	viper.SetDefault("VaultModule.KeysetConfig.KeysetName1", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetName2", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetName3", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetName4", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetValue1", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetValue2", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetValue3", "")
	viper.SetDefault("VaultModule.KeysetConfig.KeysetValue4", "")
	viper.SetDefault("VaultModule.KMSConfig.AWSKMSKey", "")
	viper.SetDefault("VaultModule.DynamoConfig.DynamoDBTableName", "")
	viper.SetDefault("VaultModule.TokenConfig.InstanceID", "")

	viper.SetDefault("VaultModule.AppDynamicsConfig.ApplicationName", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.TierName", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.NodeName", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.ControllerHostName", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.ControllerPort", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.ControllerUseSSL", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.AccountName", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.AccessKey", "")
	viper.SetDefault("VaultModule.AppDynamicsConfig.InitialTimeouts", "")
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
