package main

import (
	"log"
	"net/http"
	"pillar/internal/handlers"
)

func main() {
	new_router := handlers.NewPillarRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", new_router))
}
