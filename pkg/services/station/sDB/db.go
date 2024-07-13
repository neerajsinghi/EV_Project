package sdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/station"
	"bikeRental/pkg/services/bikeDevice/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = station.NewRepository("station")

type service struct{}

func NewService() Serv {
	return &service{}
}

func (s *service) AddStation(document entity.StationDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())

	return repo.InsertOne(document)
}

func (s *service) UpdateStation(id string, document entity.StationDB) (string, error) {
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
	update["update_at"] = primitive.NewDateTimeFromTime(time.Now())
	return repo.UpdateOne(filter, bson.M{"$set": update})
}

func (s *service) DeleteStation(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (s *service) GetStation() ([]entity.StationDB, error) {
	resp, err := repo.Find(bson.M{}, bson.M{})
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(resp); i++ {
		bikes, err := db.NewService().FindBikeByStation(resp[i].ID.Hex())
		if err == nil {
			resp[i].Stock = new(int)
			*resp[i].Stock = 0
			for j := 0; j < len(bikes); j++ {
				if bikes[j].DeviceData.Status == "online" && bikes[j].DeviceData.BatteryLevel > 20 {
					*resp[i].Stock += 1
				}
			}
		}
	}
	return resp, err
}
func (s *service) GetStationByID(id string) (entity.StationDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.FindOne(bson.M{"_id": idObject}, bson.M{})
}

func (s *service) GetNearByStation(lat, long float64, distance int) ([]entity.StationDB, error) {
	resp, err := repo.Find(bson.M{
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
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(resp); i++ {
		bikes, err := db.NewService().FindBikeByStation(resp[i].ID.Hex())
		if err == nil {
			resp[i].Stock = new(int)
			*resp[i].Stock = 0
			for j := 0; j < len(bikes); j++ {
				if bikes[j].DeviceData.Status == "online" && bikes[j].DeviceData.BatteryLevel > 20 {
					*resp[i].Stock += 1
				}
			}
		}
	}
	for i := 0; i < len(resp); i++ {
		if resp[i].Stock == nil || *resp[i].Stock == 0 {
			resp = append(resp[:i], resp[i+1:]...)
			i--
		}
	}
	return resp, err
}
