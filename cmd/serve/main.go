package main

import (
	"log"
	"net/http"
	"pillar/internal/handlers"
)

func main() {
    new_router := handlers.NewPillarRouter()
    log.Fatal(http.ListenAndServe("0.0.0.0:6969", new_router))
}


