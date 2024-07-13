package servdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/servicesRepo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = servicesRepo.NewRepository("services")

type serv struct{}

// NewService creates a new service
func NewService() IService {
	return &serv{}
}

// InsertOne inserts a new service
func (s *serv) InsertOne(document entity.ServiceDB) (string, error) {
	document.ID = primitive.NewObjectID()
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())

	return repo.InsertOne(document)
}

// FindOne finds a service by ID

func (*serv) GetAllServices() ([]entity.ServiceDB, error) {
	return repo.Find(bson.M{}, bson.M{})
}

func (s *serv) UpdateService(id string, document entity.ServiceDB) (string, error) {
	set := bson.M{}
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	if document.Name != "" {
		set["name"] = document.Name
	}
	if document.Description != "" {
		set["description"] = document.Description
	}
	if document.Price != 0 {
		set["price"] = document.Price
	}
	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return repo.UpdateOne(filter, bson.M{"$set": set})
}

func (s *serv) DeleteService(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}
