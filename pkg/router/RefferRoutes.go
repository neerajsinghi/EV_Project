package router

import reffer "bikeRental/pkg/services/reffer"

var refferRoutes = Routes{
	Route{
		"get all referral",
		"GET",
		"/referral",
		reffer.GetReferralsHandler,
	},
}
