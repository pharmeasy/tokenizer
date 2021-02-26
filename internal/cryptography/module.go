package cryptography

import (
	"context"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/errormanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/keysetmanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/tokenmanager"

	config2 "bitbucket.org/pharmaeasyteam/goframework/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/config"
)

//ModuleCrypto ...
type ModuleCrypto struct {
	config *config.TokenizerConfig
}

//New ...
func New(worldconfig config.TokenizerConfig) *ModuleCrypto {
	return &ModuleCrypto{config: &worldconfig}
}

// Init ...
func (ms *ModuleCrypto) Init(ctx context.Context, config config2.ServerConfig) {
	errormanager.SetGenericErrors()
	database.GetSession(ms.config.VaultModule.DynamoConfig.DynamoDBTableName)
	keysetmanager.LoadArnFromConfig(ms.config.VaultModule.KMSConfig.AWSKMSKey)
	keysetmanager.LoadKeysetFromConfig(map[string]string{
		ms.config.VaultModule.KeysetConfig.KeysetName1: ms.config.VaultModule.KeysetConfig.KeysetValue1,
		ms.config.VaultModule.KeysetConfig.KeysetName2: ms.config.VaultModule.KeysetConfig.KeysetValue2,
		ms.config.VaultModule.KeysetConfig.KeysetName3: ms.config.VaultModule.KeysetConfig.KeysetValue3,
		ms.config.VaultModule.KeysetConfig.KeysetName4: ms.config.VaultModule.KeysetConfig.KeysetValue4,
	})
	tokenmanager.LoadInstanceIDFromConfig(ms.config.VaultModule.TokenConfig.InstanceID)
}
