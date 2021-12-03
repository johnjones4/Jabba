package main

import (
	"log"
	"main/job"
	"net/http"
	"os"

	"github.com/johnjones4/Jabba/core"
)

func main() {
	g, err := job.NewAlertGeneratorConcrete(os.Getenv("JOB_DEFINITIONS_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	u := core.NewUpstreamConcrete(os.Getenv("UPSTREAM_URL"))

	h := initAPIServer(g, u)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), h)
	if err != nil {
		log.Fatal(err)
	}
}
