package main

import (
	"futureEVChronJobs/pkg/services/chronjobs"
	"futureEVChronJobs/pkg/services/motog"
	"log"
	"time"

	commonGo "github.com/Trestx-technology/trestx-common-go-lib"
)

func main() {

	commonGo.LoadConfig()
	go func() {
		for {
			motog.GetDataFromPullAPI()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for {
			motog.AddDeviceMoto()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for {
			chronjobs.CheckBooking()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for {
			chronjobs.CheckAndUpdateOnGoingRides()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for {
			chronjobs.GetUsersWithPlan()
			time.Sleep(time.Minute)
		}
	}()
	for {
		log.Println("Running")
		time.Sleep(time.Minute * 500)
	}

}
