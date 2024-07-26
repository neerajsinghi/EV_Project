package router

import vehicletype "bikeRental/pkg/services/vehicleType"

var vehicleTypeRoutes = Routes{
	Route{
		"Add Vehicle Type",
		"POST",
		"/vehicle/type",
		vehicletype.AddVehicleType,
	},
	Route{
		"Get All Vehicle Types",
		"GET",
		"/vehicle/type",
		vehicletype.GetAllVehicleTypes,
	},
	Route{
		"Update Vehicle Type",
		"PATCH",
		"/vehicle/type/{id}",
		vehicletype.UpdateVehicleType,
	},
	Route{
		"Delete Vehicle Type",
		"DELETE",
		"/vehicle/type/{id}",
		vehicletype.DeleteVehicleType,
	},
}
