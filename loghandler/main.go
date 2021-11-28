package main

import (
	"log"
	"main/job"
	"main/upstream"
	"net/http"
	"os"
)

func main() {
	g, err := job.NewAlertGeneratorConcrete(os.Getenv("JOB_DEFINITIONS_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	u := upstream.NewUpstreamConcrete(os.Getenv("UPSTREAM_URL"))

	h := initAPIServer(g, u)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), h)
	if err != nil {
		log.Fatal(err)
	}
}
