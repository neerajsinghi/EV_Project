package router

import account_service "bikeRental/pkg/services/account"

var accountRoutes = Routes{
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
}
