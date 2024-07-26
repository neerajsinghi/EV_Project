package router

import booking "bikeRental/pkg/services/booking"

var BookingRoutes = Routes{
	Route{
		"Add Booking",
		"POST",
		"/booking",
		booking.Book,
	},
	Route{
		"Get All Bookings",
		"GET",
		"/bookings",
		booking.GetAllBookings,
	},
	Route{
		"Get Booking All My",
		"GET",
		"/bookings/my/{id}",
		booking.GetMyAllBooking,
	},
	Route{
		"Update Booking",
		"PATCH",
		"/booking/{id}",
		booking.UpdateBooking,
	},
	Route{
		"Get Booking latest",
		"GET",
		"/bookings/my/latest/{id}",
		booking.GetMyLatestBooking,
	},
	Route{
		"Get Booking By ID",
		"GET",
		"/booking/{id}",
		booking.GetBookingByID,
	},
}
