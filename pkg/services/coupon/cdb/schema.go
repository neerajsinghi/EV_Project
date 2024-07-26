package cdb

import "bikeRental/pkg/entity"

type Coupon interface {
	AddCoupon(document entity.CouponDB) (string, error)
	UpdateCoupon(id string, document entity.CouponDB) (string, error)
	DeleteCoupon(id string) error
	GetCoupon() ([]entity.CouponDB, error)
}
