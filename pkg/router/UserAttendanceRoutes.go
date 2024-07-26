package router

import userattendance "bikeRental/pkg/services/userAttendance"

var userAttendanceRoutes = Routes{
	Route{
		"get all attendance",
		"GET",
		"/attendance",
		userattendance.GetUserAttendanceHandler,
	},
	Route{
		"get all attendance",
		"GET",
		"/attendance/{id}",
		userattendance.GetUserAttendanceByIDHandler,
	},
}
