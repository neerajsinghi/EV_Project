package router

import predefnotification "bikeRental/pkg/services/predefNotification"

var predefNotificationRoutes = Routes{
	Route{
		"get all notification templates",
		"GET",
		"/templates/notification",
		predefnotification.GetPredef,
	},
	Route{
		"add notification template",
		"POST",
		"/templates/notification",
		predefnotification.AddPredef,
	},
	Route{
		"update notification template",
		"PATCH",
		"/templates/notification",
		predefnotification.UpdatePredef,
	},
	Route{
		"delete notification template",
		"DELETE",
		"/templates/notification",
		predefnotification.DeletePredef,
	},
}
