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
	addRoutes(sub, AccountRoutes)
	addRoutes(sub, BikeDeviceRoutes)
	addRoutes(sub, BookedbikeRoutes)
	addRoutes(sub, BookingRoutes)
	addRoutes(sub, ChargerRoutes)
	addRoutes(sub, CityRoutes)
	addRoutes(sub, CouponRoutes)
	addRoutes(sub, FAQRoutes)
	addRoutes(sub, FeedbackRoutes)
	addRoutes(sub, IotBikeRoutes)
	addRoutes(sub, NotificationsRoutes)
	addRoutes(sub, PlanRoutes)
	addRoutes(sub, PredefNotificationRoutes)
	addRoutes(sub, RefferRoutes)
	addRoutes(sub, ServicesRoutes)
	addRoutes(sub, StationRoutes)
	addRoutes(sub, StatusRoutes)
	addRoutes(sub, UserAttendanceRoutes)
	addRoutes(sub, UsersRoutes)
	addRoutes(sub, VehicleTypeRoutes)
	addRoutes(sub, WalletRoutes)
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
