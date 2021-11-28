package main

import (
	"log"
	"main/store"
	"net/http"
	"os"
)

func main() {
	s, err := store.NewPGStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	h := initAPIServer(s)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), h)
	if err != nil {
		log.Fatal(err)
	}
}
