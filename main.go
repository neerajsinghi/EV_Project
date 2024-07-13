package main

import (
	"bikeRental/pkg/router"
	"bikeRental/pkg/services/chronjobs"
	utils "bikeRental/pkg/util"
	"log"
	"net/http"
	"time"

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
	go func() {
		utils.GetDataFromPullAPI()
		chronjobs.CheckBooking()
		chronjobs.CheckAndUpdateOnGoingRides()
		time.Sleep(time.Minute)
	}()
	router := router.NewRouter()
	log.Fatal(http.ListenAndServe(":1995", setupGlobalMiddleware(router)))
}
