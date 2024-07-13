package generic

import (
	"go.mongodb.org/mongo-driver/bson"
)

type BookingRepository interface {
	InsertOne(document interface{}) (string, error)
	FindOne(filter, projection bson.M) (interface{}, error)
	Find(filter, projection bson.M) ([]interface{}, error)
	UpdateOne(filter, update bson.M) (string, error)
	DeleteOne(filter bson.M) error
}
