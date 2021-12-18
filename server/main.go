package main

import (
	"encoding/json"
	"log"
	"main/status"
	"main/store"
	"net/http"
	"os"
)

func main() {
	vendorData, err := os.ReadFile(os.Getenv("VENDOR_TYPES_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	var vendorInfo map[string]string
	err = json.Unmarshal(vendorData, &vendorInfo)
	if err != nil {
		log.Fatal(err)
	}

	s, err := store.NewPGStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	se := status.NewMemoryStatusEngine(vendorInfo, s)
	go se.Start()

	h := initAPIServer(s, se)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), h)
	if err != nil {
		log.Fatal(err)
	}
}
