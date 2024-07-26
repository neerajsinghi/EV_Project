package router

import (
	account_service "bikeRental/pkg/services/account"
	"bikeRental/pkg/services/bikeDevice"
	bookedbike "bikeRental/pkg/services/bookedBike"
	"bikeRental/pkg/services/booking"
	"bikeRental/pkg/services/charger"
	"bikeRental/pkg/services/city"
	"bikeRental/pkg/services/coupon"
	faqdb "bikeRental/pkg/services/faq"
	"bikeRental/pkg/services/feedback"
	iotbike "bikeRental/pkg/services/iotBike"
	"bikeRental/pkg/services/notifications"
	"bikeRental/pkg/services/plan"
	predefnotification "bikeRental/pkg/services/predefNotification"
	"bikeRental/pkg/services/reffer"
	"bikeRental/pkg/services/services"
	"bikeRental/pkg/services/station"
	"bikeRental/pkg/services/status"
	userattendance "bikeRental/pkg/services/userAttendance"
	"bikeRental/pkg/services/users"
	vehicletype "bikeRental/pkg/services/vehicleType"
	"bikeRental/pkg/services/wallet"
	"net/http"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	Route{
		"",
		"GET",
		"/",
		hello,
	},
	Route{
		"Register",
		"POST",
		"/register",
		account_service.SignUp,
	},
	Route{
		"Forgot Password",
		"POST",
		"/forgot/password",
		account_service.ForgetPasswordOTPLink,
	},
	Route{
		"Forgot Password",
		"POST",
		"/verify/forgot/password",
		account_service.VerifyAndUpdatePassword,
	},
	Route{
		"Login",
		"POST",
		"/email/login",
		account_service.LoginUsingPassword,
	},
	Route{
		"LoginPhone",
		"POST",
		"/phone/login",
		account_service.LoginUsingPhone,
	},
	Route{
		"VerifyOTP",
		"POST",
		"/verify/otp",
		account_service.VerifyOTPAndSendToken,
	},
	Route{
		"Get ALL Bikes",
		"GET",
		"/bikes",
		iotbike.GetAll,
	},
	Route{
		"Get Nearest Bikes",
		"GET",
		"/bikes/near",
		iotbike.GetNearest,
	},
	Route{
		"Add Booking",
		"POST",
		"/booking",
		booking.Book,
	},
	Route{
		"Get All Bookings",
		"GET",
		"/bookings",
		booking.GetAllBookings,
	},
	Route{
		"Get Booking All My",
		"GET",
		"/bookings/my/{id}",
		booking.GetMyAllBooking,
	},
	Route{
		"Update Booking",
		"PATCH",
		"/booking/{id}",
		booking.UpdateBooking,
	},
	Route{
		"Get Booking latest",
		"GET",
		"/bookings/my/latest/{id}",
		booking.GetMyLatestBooking,
	},
	Route{
		"Get Booking By ID",
		"GET",
		"/booking/{id}",
		booking.GetBookingByID,
	},
	Route{
		"Add Wallet",
		"POST",
		"/wallet",
		wallet.AddWallet,
	},
	Route{
		"Get Wallet",
		"GET",
		"/wallet/{id}",
		wallet.GetMyWallet,
	},
	Route{
		"Get Wallet All",
		"GET",
		"/wallet",
		wallet.GetAllWallets,
	},
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
	Route{
		"Add Bike",
		"POST",
		"/bike",
		bikeDevice.AddBikeDevice,
	},
	Route{
		"Get All Bikes",
		"GET",
		"/bike",
		bikeDevice.GetAll,
	},
	Route{
		"Get Bike By Station",
		"GET",
		"/bike/{stationID}",
		bikeDevice.GetBikeDevicesByStation,
	},
	Route{
		"Get Bike By Station",
		"GET",
		"/bike/device/{id}",
		bikeDevice.GetBikeDevicesByDeviceID,
	},
	Route{
		"Update Bike",
		"PATCH",
		"/bike/{id}",
		bikeDevice.UpdateBikeDevice,
	},
	Route{
		"Delete Bike",
		"DELETE",
		"/bike/{id}",
		bikeDevice.DeleteBikeDevice,
	},
	Route{
		"Add Coupon",
		"POST",
		"/coupon",
		coupon.AddCoupon,
	},
	Route{
		"Get All Coupons",
		"GET",
		"/coupon",
		coupon.GetCoupons,
	},
	Route{
		"Get Coupon By ID",
		"PATCH",
		"/coupon/{id}",
		coupon.UpdateCoupon,
	},
	Route{
		"Delete Coupon",
		"DELETE",
		"/coupon/{id}",
		coupon.DeleteCoupon,
	},
	Route{
		"Add Station",
		"POST",
		"/station",
		station.AddStation,
	},
	Route{
		"Get All Stations",
		"GET",
		"/station",
		station.GetAllStations,
	},
	Route{
		"GetStation nearby",
		"GET",
		"/station/near",
		station.GetNearByStations,
	},
	Route{
		"GetStation nearby",
		"GET",
		"/station/id/{id}",
		station.GetStationsByID,
	},
	Route{
		"Update Station",
		"PATCH",
		"/station/{id}",
		station.UpdateStation,
	},
	Route{
		"Delete Station",
		"DELETE",
		"/station/{id}",
		station.DeleteStation,
	},
	Route{
		"Add Vehicle Type",
		"POST",
		"/vehicle/type",
		vehicletype.AddVehicleType,
	},
	Route{
		"Get All Vehicle Types",
		"GET",
		"/vehicle/type",
		vehicletype.GetAllVehicleTypes,
	},
	Route{
		"Update Vehicle Type",
		"PATCH",
		"/vehicle/type/{id}",
		vehicletype.UpdateVehicleType,
	},
	Route{
		"Delete Vehicle Type",
		"DELETE",
		"/vehicle/type/{id}",
		vehicletype.DeleteVehicleType,
	},
	Route{
		"Insert Service",
		"POST",
		"/service",
		services.AddService,
	},
	Route{
		"Get All Services",
		"GET",
		"/service",
		services.GetService,
	},
	Route{
		"Update Service",
		"PATCH",
		"/service/{id}",
		services.UpdateService,
	},
	Route{
		"Delete Service",
		"DELETE",
		"/service/{id}",
		services.DeleteService,
	},
	Route{
		"Add Charger",
		"POST",
		"/charger",
		charger.AddCharger,
	},
	Route{
		"Get All chargers",
		"GET",
		"/charger",
		charger.GetAllChargers,
	},
	Route{
		"GetSt nearby",
		"GET",
		"/charger/near",
		charger.GetNearByChargers,
	},
	Route{
		"Update Station",
		"PATCH",
		"/charger/{id}",
		charger.UpdateCharger,
	},
	Route{
		"Delete Station",
		"DELETE",
		"/charger/{id}",
		charger.DeleteCharger,
	},
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
	//Feedback
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
	Route{
		"statistics",
		"GET",
		"/statistics",
		status.Statistics,
	},
	Route{
		"get all cities",
		"GET",
		"/cities",
		city.GetAllCitiesHandler,
	},
	Route{
		"get city",
		"GET",
		"/city/{id}",
		city.GetCityHandler,
	},
	Route{
		"update city",
		"PATCH",
		"/city/{id}",
		city.UpdateCityHandler,
	},
	Route{
		"delete city",
		"DELETE",
		"/city/{id}",
		city.DeleteCityHandler,
	},
	Route{
		"add city",
		"POST",
		"/city",
		city.AddCityHandler,
	},
	Route{
		"get city by location",
		"GET",
		"/cities/in",
		city.InCityHandler,
	},
	Route{
		"get ongoing rides",
		"GET",
		"/rides/ongoing",
		bookedbike.GetBookedBike,
	},
	Route{
		"get all referral",
		"GET",
		"/referral",
		reffer.GetReferralsHandler,
	},
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

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
