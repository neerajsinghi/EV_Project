package router

import station "bikeRental/pkg/services/station"

var StationRoutes = Routes{
	Route{
		"Add Station",
		"POST",
		"/station",
		station.AddStation,
	},
	Route{
		"Get All Stations",
		"GET",
		"/station",
		station.GetAllStations,
	},
	Route{
		"GetStation nearby",
		"GET",
		"/station/near",
		station.GetNearByStations,
	},
	Route{
		"GetStation nearby",
		"GET",
		"/station/id/{id}",
		station.GetStationsByID,
	},
	Route{
		"Update Station",
		"PATCH",
		"/station/{id}",
		station.UpdateStation,
	},
	Route{
		"Delete Station",
		"DELETE",
		"/station/{id}",
		station.DeleteStation,
	},
}
