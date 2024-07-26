package router

import bookedbike "bikeRental/pkg/services/bookedBike"

var bookedBikeRoutes = Routes{
	Route{
		"get ongoing rides",
		"GET",
		"/rides/ongoing",
		bookedbike.GetBookedBike,
	},
}
