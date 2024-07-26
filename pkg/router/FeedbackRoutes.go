package router

import feedback "bikeRental/pkg/services/feedback"

var FeedbackRoutes = Routes{
	Route{
		"Add Feedback",
		"POST",
		"/feedback",
		feedback.AddFeedback,
	},
	Route{
		"Get All Feedbacks",
		"GET",
		"/feedback",
		feedback.GetFeedbacks,
	},
	Route{
		"Delete Feedback",
		"DELETE",
		"/feedback/{id}",
		feedback.DeleteFeedback,
	},
}
