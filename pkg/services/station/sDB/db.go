package sdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/generic"
	"bikeRental/pkg/repo/station"
	"bikeRental/pkg/services/bikeDevice/db"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = station.NewRepository("station")
var repoGeneric = generic.NewRepository("station")

type service struct{}

func NewService() Serv {
	return &service{}
}

func (s *service) AddStation(document entity.StationDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	if document.Location.Coordinates != nil {
		if len(document.Location.Coordinates) == 2 {
			if document.Location.Coordinates[0] == 0 || document.Location.Coordinates[1] == 0 {
				return "", errors.New("invalid coordinates")
			}
		} else {
			return "", errors.New("invalid coordinates")
		}
	} else {
		return "", errors.New("coordinates is required")
	}
	if document.SupervisorID == "" {
		return "", errors.New("supervisor id is required")
	} else {
		if _, err := primitive.ObjectIDFromHex(document.SupervisorID); err != nil {
			return "", errors.New("invalid supervisor id")
		}
	}
	if document.Location == nil {
		return "", errors.New("location is required")
	}
	if document.Address.City == "" {
		return "", errors.New("city is required")
	}
	if document.Name == "" {
		return "", errors.New("name is required")
	}
	data, err := repo.InsertOne(document)
	if err != nil {
		return "", errors.New("error in inserting station")
	}
	return data, nil
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
	if document.ServicesAvailable != nil {
		update["services_available"] = document.ServicesAvailable
	}
	if document.Public != nil {
		update["public"] = document.Public
	}
	if document.Stock != nil {
		update["stock"] = document.Stock
	}
	if document.SupervisorID != "" {
		if _, err := primitive.ObjectIDFromHex(document.SupervisorID); err != nil {
			return "", errors.New("invalid supervisor id")
		}
		update["supervisor_id"] = document.SupervisorID
	}
	update["update_at"] = primitive.NewDateTimeFromTime(time.Now())
	data, err := repo.UpdateOne(filter, bson.M{"$set": update})
	if err != nil {
		return "", errors.New("error in updating station")
	}
	return data, nil
}

func (s *service) DeleteStation(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (s *service) GetStation(userId, stationId string) ([]entity.StationDB, error) {
	pipeline := bson.A{}

	filter := bson.M{}
	if userId != "" {
		filter = bson.M{"supervisor_id": userId}
		pipeline = bson.A{
			bson.D{
				{Key: "$match", Value: filter},
			},
		}
	}
	if stationId != "" {
		idObject, _ := primitive.ObjectIDFromHex(stationId)
		filter["_id"] = idObject
		pipeline = bson.A{
			bson.D{
				{Key: "$match", Value: filter},
			},
		}
	}
	pipeline = append(pipeline, bson.A{
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "iotBile"},
					{Key: "localField", Value: "device_id"},
					{Key: "foreignField", Value: "deviceId"},
					{Key: "as", Value: "bikes"},
				},
			},
		},
		bson.D{{Key: "$addFields", Value: bson.D{{Key: "uID", Value: bson.D{{Key: "$toObjectId", Value: "$supervisor_id"}}}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "users"},
					{Key: "localField", Value: "uID"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "supervisor"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$supervisor"}}}},
	}...)

	cursor, err := repoGeneric.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("stations not found")
	}
	defer cursor.Close(context.Background())
	var stations []entity.StationDB
	for cursor.Next(context.Background()) {
		var station entity.StationDB
		if err = cursor.Decode(&station); err != nil {
			continue
		}
		if station.Stock == nil {
			station.Stock = new(int)
			*station.Stock = 0
		}
		for _, dev := range station.Bikes {
			if dev.Status == "online" && dev.BatteryLevel > 20 {
				*station.Stock += 1
			}
		}

		stations = append(stations, station)
	}
	if len(stations) == 0 {
		return nil, errors.New("stations not found")
	}
	return stations, err
}
func (s *service) GetStationByID(id string) (entity.StationDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	data, err := repo.FindOne(bson.M{"_id": idObject}, bson.M{})
	if err != nil {
		return entity.StationDB{}, errors.New("station not found")
	}
	return data, nil
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
		return nil, errors.New("stations not found")
	}
	for i := 0; i < len(resp); i++ {
		bikes, err := db.NewService().FindBikeByStation(resp[i].ID.Hex())
		if err == nil {
			resp[i].Stock = new(int)
			*resp[i].Stock = 0
			for j := 0; j < len(bikes); j++ {
				if bikes[j].DeviceData != nil && bikes[j].DeviceData.Status == "online" && bikes[j].DeviceData.BatteryLevel > 20 {
					*resp[i].Stock += 1
				}
			}
		}
	}
	for i := 0; i < len(resp); i++ {
		if resp[i].Stock == nil {
			resp[i].Stock = new(int)
			*resp[i].Stock = 0
		}
	}
	return resp, err
}
