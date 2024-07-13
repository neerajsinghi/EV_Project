package db

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/bikeDevice"
	bikeDB "bikeRental/pkg/services/iotBike/db"
	vdb "bikeRental/pkg/services/vehicleType/vDB"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

// FindBikeByDeviceID implements IOtBike.
func (s *service) FindBikeByDeviceID(deviceId string) ([]entity.DeviceInfo, error) {
	deviceIdIn, _ := strconv.Atoi(deviceId)
	data, err := repo.Find(bson.M{"device_id": deviceIdIn}, bson.M{})
	if err != nil {
		return nil, err
	}
	vehicleTypeIDs := make([]primitive.ObjectID, 0)
	deviceIDs := make([]int, 0)

	for _, v := range data {
		vIDObj, _ := primitive.ObjectIDFromHex(v.VehicleTypeID)
		vehicleTypeIDs = append(vehicleTypeIDs, vIDObj)
		deviceIDs = append(deviceIDs, v.DeviceID)
	}
	vehicleTypeList, err := vdb.GetVehicleType(vehicleTypeIDs)
	if err != nil {
		return nil, err
	}
	deviceList, err := bikeDB.GetBike(deviceIDs)
	if err != nil {
		return nil, err
	}
	for i, v := range data {
		for _, vt := range vehicleTypeList {
			if v.VehicleTypeID == vt.ID.Hex() {
				data[i].VehicleType = &vt
				break
			}
		}
		for _, d := range deviceList {
			if v.DeviceID == d.DeviceId {
				data[i].DeviceData = &d
				break
			}
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
	data, err := repo.Find(bson.M{}, bson.M{})
	if err != nil {
		return nil, err
	}
	vehicleTypeIDs := make([]primitive.ObjectID, 0)
	deviceIDs := make([]int, 0)
	for _, v := range data {
		vIDObj, _ := primitive.ObjectIDFromHex(v.VehicleTypeID)
		vehicleTypeIDs = append(vehicleTypeIDs, vIDObj)
		deviceIDs = append(deviceIDs, v.DeviceID)
	}
	vehicleTypeList, err := vdb.GetVehicleType(vehicleTypeIDs)
	if err != nil {
		return nil, err
	}
	deviceList, err := bikeDB.GetBike(deviceIDs)
	if err != nil {
		return nil, err
	}
	for i, v := range data {
		for _, vt := range vehicleTypeList {
			if v.VehicleTypeID == vt.ID.Hex() {
				data[i].VehicleType = &vt
				break
			}
		}
		for _, d := range deviceList {
			if v.DeviceID == d.DeviceId {
				data[i].DeviceData = &d
				break
			}
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
	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return repo.UpdateOne(bson.M{"device_id": deviceID}, bson.M{"$set": set})
}

func (*service) FindBikeByStation(stationID string) ([]entity.DeviceInfo, error) {
	data, err := repo.Find(bson.M{"station_id": stationID}, bson.M{})
	if err != nil {
		return nil, err
	}
	vehicleTypeIDs := make([]primitive.ObjectID, 0)
	deviceIDs := make([]int, 0)

	for _, v := range data {
		vIDObj, _ := primitive.ObjectIDFromHex(v.VehicleTypeID)
		vehicleTypeIDs = append(vehicleTypeIDs, vIDObj)
		deviceIDs = append(deviceIDs, v.DeviceID)
	}
	vehicleTypeList, err := vdb.GetVehicleType(vehicleTypeIDs)
	if err != nil {
		return nil, err
	}
	deviceList, err := bikeDB.GetBike(deviceIDs)
	if err != nil {
		return nil, err
	}
	for i, v := range data {
		for _, vt := range vehicleTypeList {
			if v.VehicleTypeID == vt.ID.Hex() {
				data[i].VehicleType = &vt
				break
			}
		}
		for _, d := range deviceList {
			if v.DeviceID == d.DeviceId {
				data[i].DeviceData = &d
				break
			}
		}
	}
	return data, nil
}
