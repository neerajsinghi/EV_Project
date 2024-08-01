package cdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/coupon"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = coupon.NewRepository("coupon")

type couponS struct{}

func NewCoupon() Coupon {
	return &couponS{}
}

func (c *couponS) AddCoupon(document entity.CouponDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())

	return repo.InsertOne(document)
}

func GetCouponByCode(code string) (*entity.CouponReport, error) {
	data, err := repo.FindOne(bson.M{"code": code}, bson.M{})
	if err != nil {
		return nil, err
	}
	return &data, nil
}
func GetCouponByType(couponType string) (*entity.CouponReport, error) {
	data, err := repo.FindOne(bson.M{"coupon_type": couponType}, bson.M{})
	if err != nil {
		return nil, err
	}
	return &data, nil
}
func (c *couponS) UpdateCoupon(id string, document entity.CouponDB) (string, error) {
	var updateFields bson.M
	idObject, _ := primitive.ObjectIDFromHex(id)
	document.ID = idObject
	conv, _ := bson.Marshal(document)
	bson.Unmarshal(conv, &updateFields)
	updateFields["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": updateFields})
}

func (c *couponS) DeleteCoupon(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (c *couponS) GetCoupon() ([]entity.CouponDB, error) {
	return repo.Find(bson.M{}, bson.M{})
}
