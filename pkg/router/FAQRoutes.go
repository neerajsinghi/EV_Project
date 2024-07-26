package router

import faqdb "bikeRental/pkg/services/faq"

var FAQRoutes = Routes{
	Route{
		"Get All Faqs",
		"GET",
		"/faq",
		faqdb.GetAllFaq,
	},
	Route{
		"Add Faq",
		"POST",
		"/faq",
		faqdb.AddFaq,
	},
	Route{
		"Update Faq",
		"PATCH",
		"/faq/{id}",
		faqdb.UpdateFaq,
	},
	Route{
		"Delete Faq",
		"DELETE",
		"/faq/{id}",
		faqdb.DeleteFaq,
	},
}
