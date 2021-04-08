package cryptography

import (
	"fmt"
	"net/http"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/metrics"
	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

//NewRelic newrelic middleware
func NewRelic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if metrics.GetNewRelicApp() != nil {
			logging.GetLogger().Info("Newrelic app created")
			txn := metrics.GetNewRelicApp().StartTransaction(r.URL.Path)
			logging.GetLogger().Info("Transaction started")
			w = txn.SetWebResponse(w)
			txn.SetWebRequestHTTP(r)
			defer txn.End()
			r = newrelic.RequestWithTransactionContext(r, txn)
			fmt.Println(r.Context())
			logging.GetLogger().Info("Transaction sent")
			next.ServeHTTP(w, r)
			rctx := chi.RouteContext(r.Context())
			if rctx.RoutePattern() != "" {
				txn.SetName(rctx.RoutePattern())
			}

		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)

}
