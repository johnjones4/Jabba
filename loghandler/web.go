package main

import (
	"main/job"
	"main/routes"
	"main/upstream"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/jsonschema"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/openapi"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/rest/response"
)

func getStatus(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func initAPIServer(g job.AlertGenerator, u upstream.Upstream) http.Handler {
	apiSchema := &openapi.Collector{}

	validatorFactory := jsonschema.NewFactory(apiSchema, apiSchema)
	decoderFactory := request.NewDecoderFactory()
	decoderFactory.ApplyDefaults = true
	decoderFactory.SetDecoderFunc(rest.ParamInPath, chirouter.PathToURLValues)

	r := chirouter.NewWrapper(chi.NewRouter())

	r.Use(
		middleware.Recoverer,                          // Panic recovery.
		nethttp.OpenAPIMiddleware(apiSchema),          // Documentation collector.
		request.DecoderMiddleware(decoderFactory),     // Request decoder setup.
		request.ValidatorMiddleware(validatorFactory), // Request validator setup.
		response.EncoderMiddleware,                    // Response encoder setup.
		middleware.Logger,
	)

	r.Method(http.MethodGet, "/api", http.HandlerFunc(getStatus))
	r.Method(http.MethodPost, "/api/job", nethttp.NewHandler(routes.SaveJobUsecase(g, u)))

	return r
}
