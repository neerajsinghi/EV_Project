package main

import (
	"bikeRental/pkg/router"
	"bikeRental/pkg/services/chronjobs"
	"bikeRental/pkg/services/motog"
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
		for true {
			utils.GetDataFromPullAPI()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for true {
			motog.AddDeviceMoto()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for true {
			chronjobs.CheckBooking()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for true {
			chronjobs.CheckAndUpdateOnGoingRides()
			time.Sleep(time.Minute)
		}
	}()
	router := router.NewRouter()
	log.Fatal(http.ListenAndServe(":1995", setupGlobalMiddleware(router)))
}
