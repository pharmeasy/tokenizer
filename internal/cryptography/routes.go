package cryptography

import "github.com/go-chi/chi"

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {
	router.Post("/v1/data/set", c.getTokens)
	router.Get("/v1/data/get", c.getData)

}
