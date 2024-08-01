package router

import selectedcolumns "bikeRental/pkg/services/selectedColumns"

// Routes is a slice of Route
var selectedColRoutes = Routes{
	Route{
		"Add Selected Column",
		"POST",
		"/selectedColumn",
		selectedcolumns.SelectColumnsHandler,
	},
	Route{
		"Get Selected Column",
		"GET",
		"/selectedColumn/{id}/{table}",
		selectedcolumns.GetAllSelectedColumnsHandler,
	},
	Route{
		"Get All Selected Column",
		"GET",
		"/selectedColumn/{id}",
		selectedcolumns.GetSelectedColumnsHandler,
	},
	Route{
		"Delete Selected Column",
		"DELETE",
		"/selectedColumn/{id}/{table}",
		selectedcolumns.DeleteColumnsHandler,
	},
}
