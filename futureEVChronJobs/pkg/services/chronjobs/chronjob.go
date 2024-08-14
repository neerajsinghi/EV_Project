package chronjobs

import (
	"fmt"
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/bikeDevice"
	bookedlogic "futureEVChronJobs/pkg/services/bookedBike/logic"
	bdb "futureEVChronJobs/pkg/services/booking/db"
	"futureEVChronJobs/pkg/services/city"
	"futureEVChronJobs/pkg/services/motog"
	"futureEVChronJobs/pkg/services/notifications/notify"
	pdb "futureEVChronJobs/pkg/services/plan/pDB"
	predefnotification "futureEVChronJobs/pkg/services/predefNotification"
	wdb "futureEVChronJobs/pkg/services/wallet/db"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckBooking() {
	bookings, err := bdb.GetAllHourlyBookings()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, booking := range bookings {
		planList, err := pdb.NewService().GetPlans("hourly", booking.City)
		if err != nil {
			fmt.Println(err)
			return
		}

		wallet, err := wdb.NewService().FindMy(booking.ProfileID)
		if err != nil {
			fmt.Println(err)
			return
		}
		timeSpent := 0
		walletAmount := wallet.TotalBalance
		extendedPrice := 0.0
		extendedTime := 0
		city, err := city.InCity(booking.BikeWithDevice.Location.Coordinates[1], booking.BikeWithDevice.Location.Coordinates[0])
		if err != nil || city.Name != booking.City {
			fmt.Println(err)
			if booking.Profile.FirebaseToken != nil {
				predef, err := predefnotification.Get("outOfGeofence")
				if err == nil && predef.Name == "outOfGeofence" {
					notify.NewService().SendNotification(predef.Title, predef.Body, booking.Profile.ID.Hex(), predef.Type, *booking.Profile.FirebaseToken)
				}
			}
			if booking.BikeWithDevice.Type == "moto" {
				motog.ImmoblizeDevice(1, booking.BikeWithDevice.Name)
			} else {
				motog.ImmoblizeDeviceRoadcast(booking.DeviceID, "engineStop")
			}
			filter := bson.M{"device_id": booking.DeviceID}
			repoBike := bikeDevice.NewRepository("bikeDevice")

			set := bson.M{"$set": bson.M{"immobilizeds": true}}
			repoBike.UpdateOne(filter, bson.M{"$set": set})
		}
		sort.Slice(planList, func(i, j int) bool {
			return planList[i].EndingMinutes < planList[j].EndingMinutes
		})
		for i, plan := range planList {
			if plan.EndingMinutes != 0 {
				if walletAmount == plan.Price {
					// Calculate the time this plan can provide

					// Add the time to the total timeSpent
					timeSpent = plan.EndingMinutes
					// Deduct the plan's price from the wallet

					walletAmount -= plan.Price

					break

				} else if i > 0 && (planList[i-1].Price < walletAmount && walletAmount < plan.Price) {
					timeSpent = planList[i-1].EndingMinutes
					walletAmount -= planList[i-1].Price
					break
				}
			} else if plan.EveryXMinutes != 0 {
				extendedPrice = plan.Price
				extendedTime = plan.EveryXMinutes
			}
		}
		if walletAmount > 0 {
			timeEx := int(walletAmount/extendedPrice) * extendedTime
			if timeEx >= 1 {
				timeSpent += timeEx
				walletAmount -= float64(timeSpent/extendedTime) * (extendedPrice)
			}
		}

		bdb.AddTimeRemaining(booking.ID.Hex(), timeSpent-int(time.Now().Unix()/60)+int(booking.StartTime/60))
		totalRideTime := int(time.Now().Unix()/60) - int(booking.StartTime/60)
		if totalRideTime == (timeSpent - 10) {
			if booking.Profile.FirebaseToken != nil {
				predef, err := predefnotification.Get("lastTenMinutes")
				if err == nil && predef.Name == "lastTenMinutes" {
					notify.NewService().SendNotification(predef.Title, predef.Body, booking.Profile.ID.Hex(), predef.Type, *booking.Profile.FirebaseToken)
				}
			}

		}
		if timeSpent <= totalRideTime {
			if booking.Profile.FirebaseToken != nil {
				predef, err := predefnotification.Get("timeExpired")
				if err == nil && predef.Name == "timeExpired" {
					notify.NewService().SendNotification(predef.Title, predef.Body, booking.Profile.ID.Hex(), predef.Type, *booking.Profile.FirebaseToken)
				}
			}

			//use money from wallet
			walletN := &entity.WalletS{
				ID:          primitive.NewObjectID(),
				UserID:      booking.ProfileID,
				UsedMoney:   wallet.TotalBalance - walletAmount,
				BookingID:   booking.ID.Hex(),
				Description: "Time expired for booking",
			}
			wdb.NewService().InsertOne(*walletN)
			if booking.BikeWithDevice.Type == "moto" {
				motog.ImmoblizeDevice(1, booking.BikeWithDevice.Name)
			} else {
				motog.ImmoblizeDeviceRoadcast(booking.DeviceID, "engineStop")
			}
			//stop booking
			totalDistance := booking.BikeWithDevice.TotalDistanceFloat
			bdb.ChangeStatusStopped(booking.ID.Hex(), wallet.TotalBalance-walletAmount, time.Now().Unix(), totalDistance)
		}
	}

}

func CheckAndUpdateOnGoingRides() {
	bookings, err := bdb.GetAllStartedBookings()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, booking := range bookings {
		booked := entity.BookedBikesDB{
			BookingID:   booking.ID.Hex(),
			UserID:      booking.ProfileID,
			OnGoing:     true,
			Booking:     booking,
			Coordinates: booking.BikeWithDevice.Location.Coordinates,
		}
		if booking.StartingStation != nil {
			booked.StartingStation = booking.StartingStation.Name
		}
		if booking.EndingStation != nil {
			booked.EndStation = booking.EndingStation.Name
		}
		if booking.Profile != nil {
			booked.UserName = booking.Profile.Name
		}
		if booking.BikeWithDevice != nil {
			booked.DeviceName = booking.BikeWithDevice.Name
			booked.DeviceId = booking.BikeWithDevice.DeviceId
		}
		bookedlogic.AddBookedBike(booked)
	}
}
func AddMoneyToWallet(wallet entity.WalletS) {
	wdb.NewService().InsertOne(wallet)
}
