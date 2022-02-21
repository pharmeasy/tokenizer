package cryptography

import "github.com/go-chi/chi"

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {
	router.Post("/v1/encrypt", c.encrypt)
	router.Post("/v1/decrypt", c.decrypt)
	router.Post("/v1/metadata", c.getMetaData)
	router.Put("/v1/metadata/update", c.updateMetadata)
	router.Post("/v1/hash/generate", generateHash)
}
