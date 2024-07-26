package router

import status "bikeRental/pkg/services/status"

var statusRoutes = Routes{
	Route{
		"statistics",
		"GET",
		"/statistics",
		status.Statistics,
	},
}
