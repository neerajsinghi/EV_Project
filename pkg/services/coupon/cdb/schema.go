package cdb

import "bikeRental/pkg/entity"

type Coupon interface {
	AddCoupon(document entity.CouponDB) (string, error)
	UpdateCoupon(id string, document entity.CouponDB) (string, error)
	DeleteCoupon(id string) error
	GetCoupon() ([]entity.CouponDB, error)
	GetCouponByType(couponType string) ([]entity.CouponDB, error)
	GetCouponForUser(userID, couponType, city string) ([]entity.CouponDB, error)
	GetCouponByCity(city string) ([]entity.CouponDB, error)
	GetCouponByCityAndType(city, couponType string) ([]entity.CouponDB, error)
}
