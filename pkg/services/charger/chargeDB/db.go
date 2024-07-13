package chargeDB

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/charger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = charger.NewRepository("charger")

type service struct{}

func NewService() Serv {
	return &service{}
}

func (s *service) AddCharger(document entity.ChargerDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())

	return repo.InsertOne(document)
}

func (s *service) UpdateCharger(id string, document entity.ChargerDB) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	update := bson.M{}
	if document.Status != "" {
		update["status"] = document.Status
	}
	if document.Description != "" {
		update["description"] = document.Description
	}
	if document.Active != nil {
		update["active"] = document.Active
	}
	if document.Location != nil {
		update["location"] = document.Location
	}

	if document.Public != nil {
		update["public"] = document.Public
	}
	if document.Stock != nil {
		update["stock"] = document.Stock
	}
	update["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return repo.UpdateOne(filter, bson.M{"$set": update})
}

func (s *service) DeleteCharger(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (s *service) GetCharger() ([]entity.ChargerDB, error) {
	return repo.Find(bson.M{}, bson.M{})
}
func (s *service) GetChargerByID(id string) (entity.ChargerDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.FindOne(bson.M{"_id": idObject}, bson.M{})
}

func (s *service) GetNearByCharger(lat, long float64, distance int) ([]entity.ChargerDB, error) {
	return repo.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": distance,
			},
		},
	}, bson.M{})
}
