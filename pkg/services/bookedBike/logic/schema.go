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
	filter := bson.M{"on_going": true}

	if userID != "" {
		filter["user_id"] = userID
	}
	if bookingID != "" {
		filter["booking_id"] = bookingID
	}

	pipeline := bson.A{
		// Stage 1: Match bookings based on the filter
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "userObjectId", Value: bson.D{{Key: "$toObjectId", Value: "$user_id"}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "bookingObjectId", Value: bson.D{{Key: "$toObjectId", Value: "$booking_id"}}},
		}}},
		// Stage 2: Lookup the user profile from the "users" collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},              // Collection where user profiles are stored
			{Key: "localField", Value: "userObjectId"}, // Field in bookings to match
			{Key: "foreignField", Value: "_id"},        // Field in users to match
			{Key: "as", Value: "profile"},              // Store the matched profile in this array
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "booking"},               // Collection where user profiles are stored
			{Key: "localField", Value: "bookingObjectId"}, // Field in bookings to match
			{Key: "foreignField", Value: "_id"},           // Field in users to match
			{Key: "as", Value: "booking"},                 // Store the matched profile in this array
		}}},
		// Stage 3: Unwind the "profile" array (since it's a 1-to-1 relationship)
		bson.D{{Key: "$unwind", Value: "$profile"}},
		bson.D{{Key: "$unwind", Value: "$booking"}},
		// Stage 4: Project the fields to return
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "bookedBike", Value: "$$ROOT"},
			{Key: "profile", Value: 1},
			{Key: "booking", Value: 1},
		}}},
		// Stage 5: Replace the root
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{
			{Key: "bookedBike", Value: "$bookedBike"},
			{Key: "profile", Value: "$profile"},
			{Key: "booking", Value: "$booking"},
		}}}}},
	}

	return repo.Aggregate(pipeline)
}
func ChangeOnGoing(bookingID string) error {
	filter := bson.M{"booking_id": bookingID}
	update := bson.M{"$set": bson.M{"on_going": false}}
	_, err := repo.Update(filter, update)
	return err
}
