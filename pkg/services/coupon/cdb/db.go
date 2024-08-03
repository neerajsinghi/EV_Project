package cdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/coupon"
	"errors"
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

	data, err := repo.InsertOne(document)
	if err != nil {
		return "", errors.New("error in inserting coupon")
	}
	return data, nil
}

func GetCouponByCode(code string) (*entity.CouponReport, error) {
	data, err := repo.FindOne(bson.M{"code": code}, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding coupon")
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

	data, err := repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": updateFields})
	if err != nil {
		return "", errors.New("error in updating coupon")
	}
	return data, nil
}

func (c *couponS) DeleteCoupon(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	err := repo.DeleteOne(bson.M{"_id": idObject})
	if err != nil {
		return errors.New("error in deleting coupon")
	}
	return nil
}

func (c *couponS) GetCoupon() ([]entity.CouponDB, error) {
	data, err := repo.Find(bson.M{}, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding coupon")
	}
	return data, nil
}
