package vdb

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/vehicleType"
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
	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": updateFields})
}

func (s *service) DeleteVehicleType(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (s *service) GetVehicleType() ([]entity.VehicleTypeDB, error) {
	return repo.Find(bson.M{}, bson.M{})
}

func (s *service) GetVehicleTypeByID(id string) (entity.VehicleTypeDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.FindOne(bson.M{"_id": idObject}, bson.M{})
}

func GetVehicleType(vehicleTypeIDs []primitive.ObjectID) ([]entity.VehicleTypeDB, error) {
	return repo.Find(bson.M{"_id": bson.M{"$in": vehicleTypeIDs}}, bson.M{})
}
