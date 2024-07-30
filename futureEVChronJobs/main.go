package main

import (
	"futureEVChronJobs/pkg/services/chronjobs"
	"futureEVChronJobs/pkg/services/motog"
	utils "futureEVChronJobs/pkg/util"
	"time"

	commonGo "github.com/Trestx-technology/trestx-common-go-lib"
)

func main() {

	commonGo.LoadConfig()
	go func() {
		for {
			utils.GetDataFromPullAPI()
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
}
