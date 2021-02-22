package cryptography

import (
	"context"

	config2 "bitbucket.org/pharmaeasyteam/goframework/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/config"
)

//ModuleCrypto ...
type ModuleCrypto struct {
	LoadEnvModule config.LoadEnvModule
}

//New ...
func New(worldconfig config.LoadEnvModule) *ModuleCrypto {
	return &ModuleCrypto{LoadEnvModule: worldconfig}
}

// Init ...
func (ms *ModuleCrypto) Init(ctx context.Context, config config2.ServerConfig) {

}
