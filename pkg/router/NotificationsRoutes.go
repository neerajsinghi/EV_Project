package router

import notifications "bikeRental/pkg/services/notifications"

var NotificationsRoutes = Routes{
	Route{
		"Send Notification",
		"POST",
		"/notification",
		notifications.SendNotification,
	},
	Route{
		"Send multiple Notification",
		"POST",
		"/notification/multiple",
		notifications.SendMultipleNotifications,
	},
	Route{
		"Get All Notifications",
		"GET",
		"/notification",
		notifications.GetAllNotifications,
	},
}
