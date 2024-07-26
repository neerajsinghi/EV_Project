package router

import charger "bikeRental/pkg/services/charger"

var ChargerRoutes = Routes{
	Route{
		"Add Charger",
		"POST",
		"/charger",
		charger.AddCharger,
	},
	Route{
		"Get All chargers",
		"GET",
		"/charger",
		charger.GetAllChargers,
	},
	Route{
		"GetSt nearby",
		"GET",
		"/charger/near",
		charger.GetNearByChargers,
	},
	Route{
		"Update Station",
		"PATCH",
		"/charger/{id}",
		charger.UpdateCharger,
	},
	Route{
		"Delete Station",
		"DELETE",
		"/charger/{id}",
		charger.DeleteCharger,
	},
}
