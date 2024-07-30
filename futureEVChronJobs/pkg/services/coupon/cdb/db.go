package cdb

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/coupon"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = coupon.NewRepository("coupon")

func GetCouponByCode(code string) (*entity.CouponReport, error) {
	data, err := repo.FindOne(bson.M{"code": code}, bson.M{})
	if err != nil {
		return nil, err
	}
	return &data, nil
}
