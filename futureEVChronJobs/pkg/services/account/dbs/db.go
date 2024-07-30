package db

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/profile"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	repo = profile.NewProfileRepository("users")
)

func GetUser(ids []string) ([]entity.ProfileDB, error) {
	filter := bson.A{}
	for _, id := range ids {
		idObj, _ := primitive.ObjectIDFromHex(id)
		filter = append(filter, bson.M{"_id": idObj})
	}
	return repo.Find(bson.M{"$or": filter}, bson.M{})
}

func UpdateUser(id string, user entity.ProfileDB) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	set := bson.M{}

	inc := bson.M{}
	inc["green_points"] = user.GreenPoints
	inc["carbon_saved"] = user.CarbonSaved
	inc["total_travelled"] = user.TotalTravelled
	set["update_time"] = time.Now()
	setS := bson.M{"$set": set}
	setS["$inc"] = inc

	return repo.UpdateOne(filter, setS)
}
