package cryptography

import (
	"bitbucket.org/pharmaeasyteam/tokenizer/config"
	config2 "bitbucket.org/pharmaeasyteam/goframework/config"
	"context"

)

type ModuleCrypto struct {
	config 		*config.TokenizerConfig
}

func New(worldconfig  config.TokenizerConfig) *ModuleCrypto {
	return &ModuleCrypto{config : &worldconfig}	
}


func (ms *ModuleCrypto) Init(ctx context.Context, config config2.ServerConfig) {
	
}


