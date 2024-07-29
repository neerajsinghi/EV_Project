package db

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/bikeDevice"
	"errors"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

func createPipeline(filter bson.D) bson.A {
	pipeline := bson.A{}
	if filter != nil {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}
	return append(pipeline, bson.A{
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "stationIdO", Value: bson.D{{Key: "$toObjectId", Value: "$station_id"}}},
				{Key: "vehId", Value: bson.D{{Key: "$toObjectId", Value: "$vehicle_type_id"}}},
			}},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "station"},
					{Key: "localField", Value: "stationIdO"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "stations"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "iotBike"},
					{Key: "localField", Value: "device_id"},
					{Key: "foreignField", Value: "deviceId"},
					{Key: "as", Value: "device_data"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$device_data"}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "vehicleType"},
					{Key: "localField", Value: "vehId"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "vehicle_type"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$vehicle_type"}}}},
	}...)
}

// FindBikeByDeviceID implements IOtBike.
func (s *service) FindBikeByDeviceID(deviceId string) ([]entity.DeviceInfo, error) {
	deviceIdIn, _ := strconv.Atoi(deviceId)
	pipeline := createPipeline(bson.D{{Key: "device_id", Value: deviceIdIn}})
	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	for i, v := range data {
		if len(v.Stations) > 0 {
			data[i].Station = &v.Stations[0]
		}
	}
	return data, nil

}

var (
	repo = bikeDevice.NewRepository("bikeDevice")
)

// NewService creates a new service
func NewService() IOtBike {
	return &service{}
}

func (s *service) AddBikeDevice(document entity.DeviceInfo) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	device, _ := repo.Find(bson.M{"device_id": document.DeviceID}, bson.M{})
	if len(device) > 0 {
		return "", errors.New("device already exists")
	}
	return repo.InsertOne(document)
}

func (s *service) UpdateBikeDevice(id string, document entity.DeviceInfo) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	document.ID = idObject
	set := bson.M{}

	if document.StationID != nil {
		set["station_id"] = document.StationID

	} else {
		set["station_id"] = ""
	}
	if document.Status != "" {
		set["status"] = document.Status
	}
	if document.VehicleTypeID != "" {
		set["vehicle_type_id"] = document.VehicleTypeID
	}
	if document.Description != "" {
		set["description"] = document.Description
	}
	set["update_at"] = primitive.NewDateTimeFromTime(time.Now())
	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
}

func (s *service) DeleteBikeDevice(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (s *service) FindAll() ([]entity.DeviceInfo, error) {
	pipeline := createPipeline(nil)

	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	for i, v := range data {
		if len(v.Stations) > 0 {
			data[i].Station = &v.Stations[0]
		}
	}
	return data, nil
}

func DeviceBooked(deviceID int) (string, error) {
	device, _ := repo.Find(bson.M{"device_id": deviceID}, bson.M{})
	if len(device) == 0 {
		return "", nil
	}
	set := bson.M{}
	set["station_id"] = ""
	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	return repo.UpdateOne(bson.M{"device_id": deviceID}, bson.M{"$set": bson.M{"status": "booked"}})
}

func DeviceReturned(deviceID int, stationID string) (string, error) {
	set := bson.M{}
	set["station_id"] = stationID
	set["status"] = "available"
	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return repo.UpdateOne(bson.M{"device_id": deviceID}, bson.M{"$set": set})
}

func (*service) FindBikeByStation(stationID string) ([]entity.DeviceInfo, error) {
	pipeline := createPipeline(bson.D{{Key: "station_id", Value: stationID}})
	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	for i, v := range data {
		if len(v.Stations) > 0 {
			data[i].Station = &v.Stations[0]
		}
	}
	return data, nil
}
