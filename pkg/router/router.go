package router

import "github.com/gorilla/mux"

// NewRouter builds and returns a new router from routes
func NewRouter() *mux.Router {
	// When StrictSlash == true, if the route path is "/path/", accessing "/path" will perform a redirect to the former and vice versa.
	router := mux.NewRouter().StrictSlash(true)
	router.Use(Logger)
	//Path
	sub := router.PathPrefix("/api/v1").Subrouter()
	addRoutes(sub, routes)
	addRoutes(sub, accountRoutes)
	addRoutes(sub, bikeDeviceRoutes)
	addRoutes(sub, bookedBikeRoutes)
	addRoutes(sub, bookingRoutes)
	addRoutes(sub, chargerRoutes)
	addRoutes(sub, cityRoutes)
	addRoutes(sub, couponRoutes)
	addRoutes(sub, faqRoutes)
	addRoutes(sub, feedbackRoutes)
	addRoutes(sub, iotBikeRoutes)
	addRoutes(sub, notificationsRoutes)
	addRoutes(sub, planRoutes)
	addRoutes(sub, predefNotificationRoutes)
	addRoutes(sub, refferRoutes)
	addRoutes(sub, servicesRoutes)
	addRoutes(sub, stationRoutes)
	addRoutes(sub, statusRoutes)
	addRoutes(sub, userAttendanceRoutes)
	addRoutes(sub, usersRoutes)
	addRoutes(sub, vehicleTypeRoutes)
	addRoutes(sub, walletRoutes)
	return router
}

func addRoutes(sub *mux.Router, routes Routes) {
	for _, route := range routes {
		sub.
			HandleFunc(route.Pattern, route.HandlerFunc).
			Name(route.Name).
			Methods(route.Method)
	}
}
