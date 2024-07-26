package generic

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository interface {
	InsertOne(document interface{}) (string, error)
	FindOne(filter, projection bson.M) (*mongo.Cursor, error)
	Find(filter, projection bson.M) (*mongo.Cursor, error)
	UpdateOne(filter, update bson.M) (string, error)
	DeleteOne(filter bson.M) error
	Aggregate(pipeline bson.A) (*mongo.Cursor, error)
}
