package router

import bookedbike "bikeRental/pkg/services/bookedBike"

var BookedbikeRoutes = Routes{
	Route{
		"get ongoing rides",
		"GET",
		"/rides/ongoing",
		bookedbike.GetBookedBike,
	},
}
