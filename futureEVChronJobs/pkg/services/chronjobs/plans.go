package chronjobs

import (
	"context"
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/generic"
	bdb "futureEVChronJobs/pkg/services/booking/db"
	"futureEVChronJobs/pkg/services/motog"
	"futureEVChronJobs/pkg/services/notifications/notify"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repoG = generic.NewRepository("users")

/*
plan_start_time
plan_active
plan_end_time
plan_remaining_time
*/
type profileLocal struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	FirebaseToken     string             `json:"firebase_token" bson:"firebase_token"`
	Plan              entity.PlanDB      `json:"plan" bson:"plan"`
	Booking           entity.BookingDB   `json:"bookings" bson:"bookings"`
	Bike              entity.IotBikeDB   `json:"bike" bson:"bike"`
	PlanStartTime     int64              `json:"plan_start_time" bson:"plan_start_time"`
	PlanActive        bool               `json:"plan_active" bson:"plan_active"`
	PlanEndTime       int64              `json:"plan_end_time" bson:"plan_end_time"`
	PlanRemainingTime int64              `json:"plan_remaining_time" bson:"plan_remaining_time"`
}

// get users details wo have a plan
func GetUsersWithPlan() {
	// get all users with plan
	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "plan", Value: bson.D{{Key: "$ne", Value: primitive.Null{}}}},
					{Key: "plan.type", Value: bson.D{{Key: "$ne", Value: "hourly"}}},
				},
			},
		},
		bson.D{{Key: "$addFields", Value: bson.D{{Key: "uID", Value: bson.D{{Key: "$toString", Value: "$_id"}}}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "booking"},
					{Key: "localField", Value: "uID"},
					{Key: "foreignField", Value: "profile_id"},
					{Key: "as", Value: "bookings"},
				},
			},
		},
	}
	cursor, err := repoG.Aggregate(pipeline)
	if err != nil {
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var profile profileLocal
		if err = cursor.Decode(&profile); err != nil {
			continue
		}
		// get plan start time
		if profile.PlanActive {

			planRemainingTime := (profile.PlanEndTime - time.Now().Unix())
			//convert to hours
			planRemainingTimeF := float64(planRemainingTime) / 3600.0
			if planRemainingTimeF == 1.0 || planRemainingTimeF == 1.1 {
				// send notification to user
				notify.NewService().SendNotification("Plan is about to expire", "Your plan is about to expire in 1 hour", profile.ID.Hex(), "plan Expiry", profile.FirebaseToken)
			}
			if planRemainingTime <= 0 {
				// send notification to user
				notify.NewService().SendNotification("Plan has expired", "Your plan has expired", profile.ID.Hex(), "plan Expiry", profile.FirebaseToken)
				if profile.Booking.Status == "started" {
					if profile.Bike.Type == "moto" {
						// immobilize device
						motog.ImmoblizeDevice(1, profile.Bike.Name)
					} else {
						// immobilize device
						motog.ImmoblizeDeviceRoadcast(profile.Bike.DeviceId, "engineStop")
					}
					bdb.ChangeStatusStopped(profile.Booking.ID.Hex(), 0, time.Now().Unix(), profile.Bike.TotalDistanceFloat)
				}

			}
		}
	}

}
