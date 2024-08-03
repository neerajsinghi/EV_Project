package router

import status "bikeRental/pkg/services/statistics"

var statusRoutes = Routes{
	Route{
		"statistics",
		"GET",
		"/statistics",
		status.Statistics,
	},
	Route{
		"vehicle data",
		"GET",
		"/vehicle/data/{id}",
		status.GetVehicleDataHand,
	},
	Route{
		"immobilize vehicle",
		"GET",
		"/vehicle/immobilize/{id}",
		status.ImmobilizeDevHand,
	},
}
