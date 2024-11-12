package main

import (
	"log"
	"net/http"
	"pillar/internal/router"

	"github.com/gorilla/handlers"
)

func main() {
	new_router := router.NewPillarRouter()

	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                             // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),  // Allow specific methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow specific headers
	)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", corsOptions(new_router)))
}
