package router

import coupon "bikeRental/pkg/services/coupon"

var CouponRoutes = Routes{
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
}
