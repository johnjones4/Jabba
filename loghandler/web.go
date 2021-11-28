package main

import (
	"io"
	"log"
	"main/job"
	"main/upstream"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/johnjones4/Jabba/core"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/response"
)

func getStatus(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func newSaveJobRoute(g job.AlertGenerator, u upstream.Upstream) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}

		event := core.Event{
			EventVendorType: chi.URLParam(req, "type"),
			EventVendorID:   uuid.NewString(),
			VendorInfo: map[string]string{
				"log": string(body),
			},
		}

		err = g.GenerateAlerts(&event)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}

		err = u.LogEvent(&event)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(200)
	}
}

func initAPIServer(g job.AlertGenerator, u upstream.Upstream) http.Handler {
	r := chirouter.NewWrapper(chi.NewRouter())

	r.Use(
		middleware.Recoverer,
		response.EncoderMiddleware,
		middleware.Logger,
	)

	r.Method(http.MethodGet, "/api", http.HandlerFunc(getStatus))
	r.Method(http.MethodPost, "/api/job/{type}", http.HandlerFunc(newSaveJobRoute(g, u)))

	return r
}
