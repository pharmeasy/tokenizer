package cryptography

import "github.com/go-chi/chi"

func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {
	router.Post("/v1/data/set", c.getTokens)
	router.Get("/v1/data/get", c.getData)

}

