package chronjobs

import (
	"bikeRental/pkg/entity"
	bookedlogic "bikeRental/pkg/services/bookedBike/logic"
	bdb "bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/notifications/notify"
	"bikeRental/pkg/services/users/udb"
	"bikeRental/pkg/services/wallet/db"
	wdb "bikeRental/pkg/services/wallet/db"
	"fmt"
	"time"
)

func CheckBooking() {
	bookings, err := bdb.GetAllHourlyBookings()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, booking := range bookings {
		wallet, err := wdb.NewService().FindMy(booking.ProfileID)
		if err != nil {
			fmt.Println(err)
			return
		}
		if (wallet.TotalBalance - 10) <= float64(time.Now().Unix())/60-float64(booking.StartTime)/60.0 {
			user, err := udb.NewService().GetUserById(booking.ProfileID)
			if err == nil {
				if user.FirebaseToken != nil {
					notify.NewService().SendNotification("Time Expiring", "Your booking time will Expire In 10 minutes", user.ID.Hex(), *user.FirebaseToken)
				}
			}
		}
		if (wallet.TotalBalance) <= float64(time.Now().Unix())/60-float64(booking.StartTime)/60.0 {
			//send Notification
			user, err := udb.NewService().GetUserById(booking.ProfileID)
			if err == nil {
				if user.FirebaseToken != nil {
					notify.NewService().SendNotification("Time Expired", "Your booking time is expired we will be stopping the bike", user.ID.Hex(), *user.FirebaseToken)
				}
			}
			//stop ride

			//use money from wallet
			wallet := &entity.WalletS{
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
