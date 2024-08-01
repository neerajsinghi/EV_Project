package db

import "futureEVChronJobs/pkg/entity"

type Booking interface {
	GetMyLatestBooking(userID string) (*entity.BookingOut, error)
	GetAllMyBooking(userID, bType string) ([]entity.BookingOut, error)
	GetAllBookings(status, bType, vType string) ([]entity.BookingOut, error)
}
