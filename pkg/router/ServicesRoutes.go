package router

import services "bikeRental/pkg/services/services"

var servicesRoutes = Routes{
	Route{
		"Insert Service",
		"POST",
		"/service",
		services.AddService,
	},
	Route{
		"Get All Services",
		"GET",
		"/service",
		services.GetService,
	},
	Route{
		"Update Service",
		"PATCH",
		"/service/{id}",
		services.UpdateService,
	},
	Route{
		"Delete Service",
		"DELETE",
		"/service/{id}",
		services.DeleteService,
	},
}
