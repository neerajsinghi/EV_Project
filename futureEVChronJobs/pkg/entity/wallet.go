package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletS struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	UserID          string             `json:"userId" bson:"user_id"`
	RefundableMoney float64            `json:"refundableMoney" bson:"refundable_money"`
	DepositedMoney  float64            `json:"depositedMoney" bson:"deposited_money"`
	UsedMoney       float64            `json:"usedMoney" bson:"used_money"`
	RefundedMoney   float64            `json:"refundedMoney" bson:"refunded_money"`
	PaymentID       string             `json:"paymentId" bson:"payment_id"`
	BookingID       string             `json:"bookingId" bson:"booking_id"`
	PlanID          string             `json:"planId" bson:"plan_id"`
	Plan            *PlanDB            `json:"plan" bson:"plan"`
	Status          string             `json:"status" bson:"status"`
	Booking         *BookingOut        `json:"booking" bson:"booking"`
	Description     string             `json:"description" bson:"description"`
	CreatedTime     primitive.DateTime `json:"createdTime" bson:"created_time"`
}
