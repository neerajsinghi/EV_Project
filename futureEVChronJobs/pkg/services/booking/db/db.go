package db

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/booking"
	db "futureEVChronJobs/pkg/services/account/dbs"
	"log"

	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

var (
	repo = booking.NewRepository("booking")
)

func NewService() Booking {

	return &service{}
}

// AddBooking implements Booking.

// GetAllBookings implements Booking.
func (s *service) GetAllBookings(status, bType, vType string) ([]entity.BookingOut, error) {
	filter := bson.M{}
	if status != "" && status != "all" {
		filter["status"] = status
	}
	if bType != "" && bType != "all" {
		filter["booking_type"] = bType
	}
	if vType != "" && vType != "all" {
		filter["vehicle_type"] = vType
	}
	pipeline := createPipeline(filter)

	return repo.Aggregate(pipeline)
}

func createPipeline(filter primitive.M) primitive.A {
	pipeline := bson.A{

		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "userObjectId", Value: bson.D{{Key: "$toObjectId", Value: "$profile_id"}}},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "iotBike"},
			{Key: "localField", Value: "device_id"},
			{Key: "foreignField", Value: "deviceId"},
			{Key: "as", Value: "bikeWithDevice"},
		}}},

		bson.D{{Key: "$unwind", Value: "$bikeWithDevice"}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "userObjectId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "profile"},
		}}},

		bson.D{{Key: "$unwind", Value: "$profile"}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "bookingDb", Value: "$$ROOT"},
			{Key: "bikeWithDevice", Value: 1},
			{Key: "profile", Value: 1},
		}}},

		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{
			{Key: "bookingDb", Value: "$bookingDb"},
			{Key: "bikeWithDevice", Value: "$bikeWithDevice"},
			{Key: "profile", Value: "$profile"},
		}}}}},
	}
	return pipeline
}
func GetAllHourlyBookings() ([]entity.BookingOut, error) {
	filter := bson.M{}
	filter["status"] = "started"
	filter["booking_type"] = "hourly"
	pipeline := createPipeline(filter)

	return repo.Aggregate(pipeline)
}
func GetAllStartedBookings() ([]entity.BookingOut, error) {
	filter := bson.M{}
	filter["status"] = "started"
	pipeline := createPipeline(filter)

	return repo.Aggregate(pipeline)
}
func GetMyBookingCount(userID string) (int64, error) {
	filter := bson.M{"profile_id": userID}
	return repo.Count(filter)
}

// GetAllMyBooking implements Booking.
func (s *service) GetAllMyBooking(userID, bType string) ([]entity.BookingOut, error) {
	filter := bson.M{"profile_id": userID}
	if bType != "" && bType != "all" {
		filter["status"] = bType
	}
	pipeline := createPipeline(filter)

	return repo.Aggregate(pipeline)

}

// GetMyBooking implements Booking.
func (s *service) GetMyLatestBooking(userID string) (*entity.BookingOut, error) {
	filter := bson.M{"profile_id": userID, "status": bson.M{"$ne": "completed"}}
	pipeline := createPipeline(filter)

	booking, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	if len(booking) == 0 {
		return nil, errors.New("booking not found")
	}
	if len(booking) > 1 {
		for i := 1; i < len(booking); i++ {
			if booking[i].Status == "started" {
				return &booking[i], nil
			}
		}
	}
	return &booking[0], nil
}
func (*service) GetBookingByID(id string) (*entity.BookingOut, error) {
	return GetBooking(id)
}

func GetBooking(id string) (*entity.BookingOut, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	pipeline := createPipeline(filter)

	booking, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	if len(booking) == 0 {
		return nil, errors.New("booking not found")
	}
	return &booking[0], nil
}

func AddTimeRemaining(id string, timeRemaining int) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	set := bson.M{}
	set["ride_time_remaining"] = timeRemaining
	set["update_time"] = time.Now()
	repo.UpdateOne(filter, bson.M{"$set": set})
}

func ChangeStatusStopped(id string, price float64, endTime int64, endKm float64) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	booking, _ := GetBooking(id)
	set := bson.M{}
	set["status"] = "stopped"
	set["price"] = price
	set["end_time"] = endTime
	set["end_km"] = endKm
	set["total_distance"] = (endKm - booking.StartKM)
	userTotalDist := (endKm - booking.StartKM)
	set["updated_time"] = primitive.NewDateTimeFromTime(time.Now())
	greenPoints := int64(userTotalDist * 5)
	carbonSaved := userTotalDist * 80
	profile := entity.ProfileDB{
		GreenPoints:    greenPoints,
		CarbonSaved:    carbonSaved,
		TotalTravelled: userTotalDist,
	}
	set["green_points"] = greenPoints
	set["carbon_saved"] = carbonSaved
	db.UpdateUser(booking.ProfileID, profile)
	_, err := repo.UpdateOne(filter, bson.M{"$set": set})
	log.Println("error", err)
}
