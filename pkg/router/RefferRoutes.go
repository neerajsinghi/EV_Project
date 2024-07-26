package router

import reffer "bikeRental/pkg/services/reffer"

var RefferRoutes = Routes{
	Route{
		"get all referral",
		"GET",
		"/referral",
		reffer.GetReferralsHandler,
	},
}
