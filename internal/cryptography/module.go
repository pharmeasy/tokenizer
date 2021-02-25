package cryptography

import (
	"context"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/errormanager"

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
	database.GetSession()
}
