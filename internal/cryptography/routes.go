package cryptography

import "github.com/go-chi/chi"

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {
	router.Post("/v1/encrypt", c.getTokens)
	router.Get("/v1/decrypt", c.getData)
	router.Patch("/v1/metadata/update", c.updateMeta)
}
