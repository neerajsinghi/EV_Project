package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponDB struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	ServiceType    []string           `json:"serviceType" bson:"service_type"`
	City           []string           `json:"city" bson:"city"`
	VehicleType    []string           `json:"vehicleType" bson:"vehicle_type"`
	Code           string             `json:"code" bson:"code"`
	CouponType     string             `json:"couponType" bson:"coupon_type"`
	MinValue       float64            `json:"minValue" bson:"min_value"`
	MaxValue       float64            `json:"maxValue" bson:"max_value"`
	MaxUsageByUser int                `json:"maxUsageByUser" bson:"max_usage_by_user"`
	Discount       float64            `json:"discount" bson:"discount"`
	ValidityFrom   time.Time          `json:"validFrom" bson:"valid_from"`
	ValidTill      time.Time          `json:"validTill" bson:"valid_till"`
	Description    string             `json:"description" bson:"description"`
	CreatedTime    primitive.DateTime `json:"createdTime" bson:"created_time"`
}

type CouponReport struct {
	CouponDB       `json:",inline"`
	BookingHistory []BookingDB `json:"bookingHistory" bson:"booking_history"`
	AmountSaved    float64     `json:"amountSaved" bson:"amount_saved"`
	AmountEarned   float64     `json:"amountEarned" bson:"amount_earned"`
}

/*
Coupons Tab: Should be in service type wise, the coupon should include (Code, City (Multiple), Service Type (Multiple), Coupon history, Vehicle type (Multiple), Validity (To and From), Coupon Type: Flat amount or percentage to upto XX,  minimum value), number of times it can be used by a single user. (Coupon Report which included booking history and all details associated to it and also number of times this code is used, amount saved, and amount earned), Birthday Coupon

*/
