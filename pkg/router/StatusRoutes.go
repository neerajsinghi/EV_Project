package router

import status "bikeRental/pkg/services/status"

var StatusRoutes = Routes{
	Route{
		"statistics",
		"GET",
		"/statistics",
		status.Statistics,
	},
}
