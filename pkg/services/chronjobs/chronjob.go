package chronjobs

import (
	"bikeRental/pkg/entity"
	bookedlogic "bikeRental/pkg/services/bookedBike/logic"
	bdb "bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/city"
	"bikeRental/pkg/services/motog"
	"bikeRental/pkg/services/notifications/notify"
	pdb "bikeRental/pkg/services/plan/pDB"
	predefnotification "bikeRental/pkg/services/predefNotification"
	wdb "bikeRental/pkg/services/wallet/db"
	"fmt"
	"sort"
	"time"

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
			}
		}
		sort.Slice(planList, func(i, j int) bool {
			return planList[i].EndingMinutes < planList[j].EndingMinutes
		})
		for _, plan := range planList {
			if plan.EndingMinutes != 0 {
				if walletAmount <= plan.Price {
					// Calculate the time this plan can provide

					// Add the time to the total timeSpent
					timeSpent = plan.EndingMinutes
					// Deduct the plan's price from the wallet

					walletAmount -= plan.Price

					// If there's no money left, stop iterating
					if walletAmount <= 0 {
						break
					}
				}
			} else {
				extendedPrice = plan.Price
				extendedTime = plan.EveryXMinutes
			}
		}
		if walletAmount > 0 {
			timeSpent += int(walletAmount/extendedPrice) * extendedTime
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
				UsedMoney:   wallet.TotalBalance,
				BookingID:   booking.ID.Hex(),
				Description: "Time expired for booking",
			}
			wdb.NewService().InsertOne(*walletN)
			if booking.BikeWithDevice.Type == "moto" {
				motog.ImmoblizeDevice(1, booking.BikeWithDevice.Name)
			}
			//stop booking
			bdb.ChangeStatusStopped(booking.ID.Hex(), wallet.TotalBalance, time.Now().Unix())
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
			BookingID: booking.ID.Hex(),
			UserID:    booking.ProfileID,
			OnGoing:   true,
			Booking:   booking,
		}
		bookedlogic.AddBookedBike(booked)
	}
}
func AddMoneyToWallet(wallet entity.WalletS) {
	wdb.NewService().InsertOne(wallet)
}
