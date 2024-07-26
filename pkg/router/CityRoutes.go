package router

import city "bikeRental/pkg/services/city"

var CityRoutes = Routes{
	Route{
		"get all cities",
		"GET",
		"/cities",
		city.GetAllCitiesHandler,
	},
	Route{
		"get city",
		"GET",
		"/city/{id}",
		city.GetCityHandler,
	},
	Route{
		"update city",
		"PATCH",
		"/city/{id}",
		city.UpdateCityHandler,
	},
	Route{
		"delete city",
		"DELETE",
		"/city/{id}",
		city.DeleteCityHandler,
	},
	Route{
		"add city",
		"POST",
		"/city",
		city.AddCityHandler,
	},
	Route{
		"get city by location",
		"GET",
		"/cities/in",
		city.InCityHandler,
	},
}
