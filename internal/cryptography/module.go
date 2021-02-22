package cryptography

import (
	"context"

	config2 "bitbucket.org/pharmaeasyteam/goframework/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/config"
)

//ModuleCrypto ...
type ModuleCrypto struct {
	LoadModule config.LoadModule
}

//New ...
func New(worldconfig config.LoadModule) *ModuleCrypto {
	return &ModuleCrypto{LoadModule: worldconfig}
}

// Init ...
func (ms *ModuleCrypto) Init(ctx context.Context, config config2.ServerConfig) {

}
