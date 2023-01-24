package cryptography

import (
	"github.com/go-chi/chi"
	instana "github.com/instana/go-sensor"
)

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {

	router.Post("/v1/encrypt", instana.TracingHandlerFunc(c.InstanaSensor, "test-random", c.encrypt))
	router.Post("/v1/decrypt", instana.TracingHandlerFunc(c.InstanaSensor, "test-random-2", c.decrypt))
	router.Post("/v1/metadata", instana.TracingHandlerFunc(c.InstanaSensor, "test-random-3", c.getMetaData))
	router.Put("/v1/metadata/update", instana.TracingHandlerFunc(c.InstanaSensor, "test-random-4", c.updateMetadata))
}
