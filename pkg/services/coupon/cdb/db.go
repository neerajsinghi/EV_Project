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
		coupon, _ := GetCouponByType("Referral")
		if coupon.Code != "" {
			return "", errors.New("referral coupon already exists")
		}
	}
	if strings.EqualFold(document.CouponType, "firstRide") {
		coupon, _ := repo.FindOne(bson.M{"coupon_type": "firstRide", "city": document.City}, bson.M{})
		if coupon.Code != "" {
			return "", errors.New("first ride coupon already exists")
		}
		if document.City == nil {
			return "", errors.New("city is required")
		}
		if len(document.City) == 0 {
			return "", errors.New("city is required")
		}
		if len(document.ServiceType) != 1 {
			return "", errors.New("service type is required")
		} else {
			if document.ServiceType[0] != "hourly" {
				return "", errors.New("service type can only be Ride now")
			}
		}

	}
	if strings.EqualFold(document.CouponType, "Discount") {
		if len(document.City) == 0 {
			return "", errors.New("city is required")
		}
		if document.Discount == 0 {
			return "", errors.New("discount is required")
		}
		if document.MaxValue == 0 {
			return "", errors.New("max value is required")
		}
		if document.MaxUsageByUser == 0 {
			return "", errors.New("max usage by user is required")
		}
		if len(document.ServiceType) == 0 {
			return "", errors.New("service type is required")
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
		bson.D{{Key: "$match", Value: bson.D{{Key: "service_type", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: typeC, Options: "i"}}}}, {Key: "city", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: city, Options: "i"}}}}}}},
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
		bson.D{{Key: "$match", Value: bson.D{{Key: "city", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: city, Options: "i"}}}}}}},
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
		bson.D{{Key: "$match", Value: bson.D{{Key: "service_type", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: couponType, Options: "i"}}}}}}},
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

func (c *couponS) GetCouponForUser(userID, couponType, city string) ([]entity.CouponDB, error) {
	cursor, err := repoG.Aggregate(bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "service_type", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: couponType, Options: "i"}}}}, {Key: "city", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: city, Options: "i"}}}}}}},
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
		bookings := 0
		*coupon.WalletCount = 0
		for _, booking := range coupon.Bookings {
			if booking.ProfileID == userID {
				bookings++
			}
		}
		*coupon.BookingCount = bookings
		for _, wallet := range coupon.Wallets {
			if wallet.UserID == userID {
				*coupon.WalletCount++
			}
		}
		data = append(data, coupon)
	}
	return data, nil
}
