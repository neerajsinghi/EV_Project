package router

import users "bikeRental/pkg/services/users"

var usersRoutes = Routes{
	Route{
		"Get All Users",
		"GET",
		"/users",
		users.GetUsers,
	},
	Route{
		"Get User By ID",
		"GET",
		"/users/{id}",
		users.GetUserById,
	},
	Route{
		"Update User",
		"PATCH",
		"/users/{id}",
		users.UpdateUser,
	},
	Route{
		"Delete User",
		"DELETE",
		"/users/{id}",
		users.DeleteUser,
	},
	Route{
		"Delete User",
		"DELETE",
		"/users/permanent/{id}",
		users.DeleteUserPermanently,
	},
	Route{
		"Remove Plan",
		"PATCH",
		"/users/{id}/plan/{plan_id}",
		users.RemovePlan,
	},
}
