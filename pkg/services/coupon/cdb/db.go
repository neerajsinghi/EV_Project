package cdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/coupon"
	"bikeRental/pkg/repo/generic"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = coupon.NewRepository("coupon")
var repoG = generic.NewRepository("coupon")

type couponS struct{}

func NewCoupon() Coupon {
	return &couponS{}
}

func (c *couponS) AddCoupon(document entity.CouponDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	if strings.EqualFold(document.CouponType, "Referral") {
		coupon, err := GetCouponByType("Referral")
		if err == nil || coupon != nil {
			return "", errors.New("referral coupon already exists")
		}
		document.Code = "REF" + primitive.NewObjectID().Hex()
	}
	if strings.EqualFold(document.CouponType, "firstRide") {
		coupon, err := repo.FindOne(bson.M{"coupon_type": "firstRide", "city": document.City}, bson.M{})
		if err == nil || coupon.Code != "" {
			return "", errors.New("first ride coupon already exists")
		}
		if document.City == nil {
			return "", errors.New("city is required")
		}
	}
	if strings.EqualFold(document.CouponType, "Discount") {
		if document.City == nil {
			return "", errors.New("city is required")
		}
	}

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

	cursor, err := repoG.Aggregate(bson.A{
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "booking"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "bookings"},
			}},
		},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "wallet"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "wallets"},
			},
		}},
	})
	if err != nil {
		return nil, errors.New("error in finding coupon")
	}
	defer cursor.Close(context.Background())
	var data []entity.CouponDB
	for cursor.Next(context.Background()) {
		var coupon entity.CouponDB
		if err = cursor.Decode(&coupon); err != nil {
			continue
		}
		coupon.BookingCount = new(int)
		coupon.WalletCount = new(int)
		*coupon.BookingCount = len(coupon.Bookings)
		*coupon.WalletCount = len(coupon.Wallets)
		data = append(data, coupon)
	}
	return data, nil
}
func (c *couponS) GetCouponByCityAndType(city, typeC string) ([]entity.CouponDB, error) {

	cursor, err := repoG.Aggregate(bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "service_type", Value: typeC}, {Key: "city", Value: city}}}},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "booking"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "bookings"},
			}},
		},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "wallet"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "wallets"},
			},
		}},
	})
	if err != nil {
		return nil, errors.New("error in finding coupon")
	}
	defer cursor.Close(context.Background())
	var data []entity.CouponDB
	for cursor.Next(context.Background()) {
		var coupon entity.CouponDB
		if err = cursor.Decode(&coupon); err != nil {
			continue
		}
		coupon.BookingCount = new(int)
		coupon.WalletCount = new(int)
		*coupon.BookingCount = len(coupon.Bookings)
		*coupon.WalletCount = len(coupon.Wallets)
		data = append(data, coupon)
	}
	return data, nil
}

func (c *couponS) GetCouponByCity(city string) ([]entity.CouponDB, error) {

	cursor, err := repoG.Aggregate(bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "city", Value: city}}}},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "booking"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "bookings"},
			}},
		},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "wallet"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "wallets"},
			},
		}},
	})
	if err != nil {
		return nil, errors.New("error in finding coupon")
	}
	defer cursor.Close(context.Background())
	var data []entity.CouponDB
	for cursor.Next(context.Background()) {
		var coupon entity.CouponDB
		if err = cursor.Decode(&coupon); err != nil {
			continue
		}
		coupon.BookingCount = new(int)
		coupon.WalletCount = new(int)
		*coupon.BookingCount = len(coupon.Bookings)
		*coupon.WalletCount = len(coupon.Wallets)
		data = append(data, coupon)
	}
	return data, nil
}

func (c *couponS) GetCouponByType(couponType string) ([]entity.CouponDB, error) {
	cursor, err := repoG.Aggregate(bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "service_type", Value: couponType}}}},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "booking"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "bookings"},
			}},
		},
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "wallet"},
				{Key: "localField", Value: "code"},
				{Key: "foreignField", Value: "coupon_code"},
				{Key: "as", Value: "wallets"},
			},
		}},
	})
	if err != nil {
		return nil, errors.New("error in finding coupon")
	}
	defer cursor.Close(context.Background())
	var data []entity.CouponDB
	for cursor.Next(context.Background()) {
		var coupon entity.CouponDB
		if err = cursor.Decode(&coupon); err != nil {
			continue
		}
		coupon.BookingCount = new(int)
		coupon.WalletCount = new(int)
		*coupon.BookingCount = len(coupon.Bookings)
		*coupon.WalletCount = len(coupon.Wallets)
		data = append(data, coupon)
	}
	return data, nil
}
