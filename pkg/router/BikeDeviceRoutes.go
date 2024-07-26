package router

import bikeDevice  "bikeRental/pkg/services/bikeDevice"

var bikeDeviceRoutes = Routes{
	Route{
		"Add Bike",
		"POST",
		"/bike",
		bikeDevice.AddBikeDevice,
	},
	Route{
		"Get All Bikes",
		"GET",
		"/bike",
		bikeDevice.GetAll,
	},
	Route{
		"Get Bike By Station",
		"GET",
		"/bike/{stationID}",
		bikeDevice.GetBikeDevicesByStation,
	},
	Route{
		"Get Bike By Station",
		"GET",
		"/bike/device/{id}",
		bikeDevice.GetBikeDevicesByDeviceID,
	},
	Route{
		"Update Bike",
		"PATCH",
		"/bike/{id}",
		bikeDevice.UpdateBikeDevice,
	},
	Route{
		"Delete Bike",
		"DELETE",
		"/bike/{id}",
		bikeDevice.DeleteBikeDevice,
	},
}
