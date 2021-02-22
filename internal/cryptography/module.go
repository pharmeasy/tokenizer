package cryptography

import (
	"context"

	config2 "bitbucket.org/pharmaeasyteam/goframework/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/validator"
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
	validator.SetGenericError()
}
