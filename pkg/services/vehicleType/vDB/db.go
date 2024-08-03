package vdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/vehicleType"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	repo = vehicleType.NewRepository("vehicleType")
)

type service struct{}

func NewService() Serv {
	return &service{}
}

func (s *service) AddVehicleType(document entity.VehicleTypeDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())

	return repo.InsertOne(document)
}

func (s *service) UpdateVehicleType(id string, document entity.VehicleTypeDB) (string, error) {
	var updateFields bson.M
	idObject, _ := primitive.ObjectIDFromHex(id)
	document.ID = idObject
	conv, _ := bson.Marshal(document)
	bson.Unmarshal(conv, &updateFields)
	data, err := repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": updateFields})
	if err != nil {
		return "", errors.New("error in updating vehicle type")
	}
	return data, err
}

func (s *service) DeleteVehicleType(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	err := repo.DeleteOne(bson.M{"_id": idObject})
	if err != nil {
		return errors.New("error in deleting vehicle type")
	}
	return err
}

func (s *service) GetVehicleType() ([]entity.VehicleTypeDB, error) {
	data, err := repo.Find(bson.M{}, bson.M{})
	if err != nil {
		return nil, errors.New("error in getting vehicle type")
	}
	return data, err
}

func (s *service) GetVehicleTypeByID(id string) (entity.VehicleTypeDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	data, err := repo.FindOne(bson.M{"_id": idObject}, bson.M{})
	if err != nil {
		return entity.VehicleTypeDB{}, errors.New("error in getting vehicle type")
	}
	return data, err
}

func GetVehicleType(vehicleTypeIDs []primitive.ObjectID) ([]entity.VehicleTypeDB, error) {
	data, err := repo.Find(bson.M{"_id": bson.M{"$in": vehicleTypeIDs}}, bson.M{})
	if err != nil {
		return nil, errors.New("error in getting vehicle type")
	}
	return data, err
}
