package main

import (
	"log"
	"main/core"
	"main/pollers"
	"os"
	"time"

	jabbacore "github.com/johnjones4/Jabba/core"
)

func main() {
	pollers := []core.Poller{
		pollers.NewINetPoller(),
		pollers.NewAbodePoller(os.Getenv("ABODE_USERNAME"), os.Getenv("ABODE_PASSWORD")),
		pollers.NewSitesPoller(os.Getenv("SITES_CONFIG")),
	}
	for _, p := range pollers {
		err := p.Setup()
		if err != nil {
			log.Fatal(err)
		}
	}
	u := jabbacore.NewUpstreamConcrete(os.Getenv("UPSTREAM_URL"))
	for _, p := range pollers {
		go p.Poll(u)
	}
	for {
		time.Sleep(time.Hour)
	}
}
