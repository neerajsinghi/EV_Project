package main

import (
	"bikeRental/pkg/router"
	"log"
	"net/http"

	commonGo "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/rs/cors"
)

// setupGlobalMiddleware will setup CORS
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.AllowAll().Handler
	return handleCORS(handler)
}

func main() {

	commonGo.LoadConfig()
	router := router.NewRouter()
	log.Fatal(http.ListenAndServe(":1995", setupGlobalMiddleware(router)))
}
