package chronjobs

import (
	"bikeRental/pkg/entity"
	bookedlogic "bikeRental/pkg/services/bookedBike/logic"
	bdb "bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/city"
	"bikeRental/pkg/services/notifications/notify"
	pdb "bikeRental/pkg/services/plan/pDB"
	"bikeRental/pkg/services/wallet/db"
	wdb "bikeRental/pkg/services/wallet/db"
	"fmt"
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
				notify.NewService().SendNotification("Left premises", "Your bike is going out of premise", booking.Profile.ID.Hex(), "leavingPremise", *booking.Profile.FirebaseToken)
			}
		}
		for _, plan := range planList {
			if plan.EndingMinutes != 0 {
				if walletAmount <= plan.Price {
					// Calculate the time this plan can provide

					// Add the time to the total timeSpent
					timeSpent = plan.EndingMinutes
					// Deduct the plan's price from the wallet
					walletAmount -= plan.Price

					// If there's no money left, stop iterating
					if wallet.TotalBalance <= 0 {
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
		if int(time.Now().Unix()/60)-int(booking.StartTime/60) <= (timeSpent-10) && int(time.Now().Unix()/60)-int(booking.StartTime/60) >= timeSpent-9 {
			if booking.Profile.FirebaseToken != nil {
				notify.NewService().SendNotification("Time Expiring", "Your booking time will Expire In 10 minutes Please recharge your wallet", booking.Profile.ID.Hex(), "timeExpiring", *booking.Profile.FirebaseToken)
			}

		}
		if timeSpent <= int(time.Now().Unix()/60)-int(booking.StartTime/60) {
			if booking.Profile.FirebaseToken != nil {
				notify.NewService().SendNotification("Time Expired", "Your booking time expired we will be stopping the bike", booking.Profile.ID.Hex(), "timeExpiring", *booking.Profile.FirebaseToken)
			}

			//use money from wallet
			wallet := &entity.WalletS{
				ID:          primitive.NewObjectID(),
				UserID:      booking.ProfileID,
				UsedMoney:   wallet.TotalBalance,
				BookingID:   booking.ID.Hex(),
				Description: "Time expired for booking",
			}
			db.NewService().InsertOne(*wallet)
			//stop booking
			booking := entity.BookingDB{
				Status: "stopped",
			}
			bdb.NewService().UpdateBooking(booking.ID.Hex(), booking)
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
			Bike:      *booking.BikeWithDevice,
			OnGoing:   true,
		}
		bookedlogic.AddBookedBike(booked)
	}
}
func AddMoneyToWallet(wallet entity.WalletS) {
	db.NewService().InsertOne(wallet)
}
