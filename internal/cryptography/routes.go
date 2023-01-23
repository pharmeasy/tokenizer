package cryptography

import (
	"github.com/go-chi/chi"
	instana "github.com/instana/go-sensor"
)

// RegisterRoutes ...
func (c *ModuleCrypto) RegisterRoutes(router chi.Router) {

	router.Post("/v1/encrypt", instana.TracingHandlerFunc(c.InstanaSensor, "", c.encrypt))
	router.Post("/v1/decrypt", instana.TracingHandlerFunc(c.InstanaSensor, "", c.decrypt))
	router.Post("/v1/metadata", instana.TracingHandlerFunc(c.InstanaSensor, "", c.getMetaData))
	router.Put("/v1/metadata/update", instana.TracingHandlerFunc(c.InstanaSensor, "", c.updateMetadata))
}
