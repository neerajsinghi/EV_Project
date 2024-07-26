package bookedlogic

import (
	"bikeRental/pkg/entity"
	booked "bikeRental/pkg/repo/bookedBikes"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = booked.NewRepository("bookedBikes")

func AddBookedBike(document entity.BookedBikesDB) (string, error) {
	return repo.InsertOne(document)
}

func GetBookedBike(userID, bookingID string) ([]entity.BookedBikesDB, error) {
	filter := bson.M{}

	if userID != "" {
		filter["user_id"] = userID
	}
	if bookingID != "" {
		filter["booking_id"] = bookingID
	}

	return repo.Find(filter, bson.M{})
}
func ChangeOnGoing(bookingID string) error {
	filter := bson.M{"booking_id": bookingID}
	update := bson.M{"$set": bson.M{"on_going": false}}
	_, err := repo.Update(filter, update)
	return err
}
