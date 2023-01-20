package cryptography

import (
	"github.com/go-chi/chi"
	instana "github.com/instana/go-sensor"
)

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {

	router.Post("/v1/encrypt", instana.TracingNamedHandlerFunc(c.InstanaSensor, "v1_encrypt", "/v1/encrypt", c.encrypt))
	router.Post("/v1/decrypt", instana.TracingNamedHandlerFunc(c.InstanaSensor, "v1_decrypt", "/v1/decrypt", c.decrypt))
	router.Post("/v1/metadata", instana.TracingNamedHandlerFunc(c.InstanaSensor, "v1_metadata", "/v1/metadata", c.getMetaData))
	router.Put("/v1/metadata/update", instana.TracingNamedHandlerFunc(c.InstanaSensor, "v1_metadata_update", "/v1/metadata/update", c.updateMetadata))
}
