package booking

import (
	"futureEVChronJobs/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type BookingRepository interface {
	InsertOne(document interface{}) (string, error)
	FindOne(filter, projection bson.M) (entity.BookingDB, error)
	Find(filter, projection bson.M) ([]entity.BookingDB, error)
	UpdateOne(filter, update bson.M) (string, error)
	DeleteOne(filter bson.M) error
	Count(filter bson.M) (int64, error)
	Aggregate(pipeline bson.A) ([]entity.BookingOut, error)
}
