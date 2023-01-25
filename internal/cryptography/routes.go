package cryptography

import (
	"github.com/go-chi/chi"
	instana "github.com/instana/go-sensor"
)

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {

	router.Post("/v1/encrypt", instana.TracingHandlerFunc(c.InstanaSensor, "v1-encrypt", c.encrypt))
	router.Post("/v1/decrypt", instana.TracingHandlerFunc(c.InstanaSensor, "v1-decrypt", c.decrypt))
	router.Post("/v1/metadata", instana.TracingHandlerFunc(c.InstanaSensor, "v1-get-metadata", c.getMetaData))
	router.Put("/v1/metadata/update", instana.TracingHandlerFunc(c.InstanaSensor, "v1-update-metadata", c.updateMetadata))
}
