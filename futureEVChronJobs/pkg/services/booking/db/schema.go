package db

import "futureEVChronJobs/pkg/entity"

type Booking interface {
	AddBooking(document entity.BookingDB) (string, error)
	UpdateBooking(id string, document entity.BookingDB) (string, error)
	DeleteBooking(id string) error
	GetMyLatestBooking(userID string) (*entity.BookingOut, error)
	GetAllMyBooking(userID, bType string) ([]entity.BookingOut, error)
	GetAllBookings(status, bType, vType string) ([]entity.BookingOut, error)
}
