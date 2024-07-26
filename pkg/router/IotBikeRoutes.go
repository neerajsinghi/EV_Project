package router

import iotbike "bikeRental/pkg/services/iotBike"

var IotBikeRoutes = Routes{
	Route{
		"Get ALL Bikes",
		"GET",
		"/bikes",
		iotbike.GetAll,
	},
	Route{
		"Get Nearest Bikes",
		"GET",
		"/bikes/near",
		iotbike.GetNearest,
	},
}
