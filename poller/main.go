package main

import (
	"log"
	"main/core"
	"main/pollers"
	"os"

	jabbacore "github.com/johnjones4/Jabba/core"
)

func main() {
	a := pollers.NewAbodePoller(os.Getenv("ABODE_USERNAME"), os.Getenv("ABODE_PASSWORD"))
	err := a.Authorize()
	if err != nil {
		log.Fatal(err)
	}

	u := jabbacore.NewUpstreamConcrete(os.Getenv("UPSTREAM_URL"))

	eChan := make(chan error, 255)
	pw := core.NewPollWatcher()
	a.Poll(&pw, eChan, u)
}
