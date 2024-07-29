package router

import plan "bikeRental/pkg/services/plan"

var planRoutes = Routes{
	Route{
		"Add Plan",
		"POST",
		"/plan",
		plan.AddPlan,
	},
	Route{
		"Get All Plans",
		"GET",
		"/plan",
		plan.GetAllPlans,
	},
	Route{
		"Get Plan By ID",
		"GET",
		"/plan/{id}",
		plan.GetPlanById,
	},
	Route{
		"Update Plan",
		"PATCH",
		"/plan/{id}",
		plan.UpdatePlan,
	},
	Route{
		"Delete Plan",
		"DELETE",
		"/plan/{id}",
		plan.DeletePlan,
	},
	Route{
		"Get Deposit",
		"GET",
		"/deposit",
		plan.GetDeposit,
	},
}
