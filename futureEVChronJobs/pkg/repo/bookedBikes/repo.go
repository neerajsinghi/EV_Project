package booked

import (
	"futureEVChronJobs/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	InsertOne(document interface{}) (string, error)
	FindOne(filter, projection bson.M) (entity.BookedBikesDB, error)
	Find(filter, projection bson.M) ([]entity.BookedBikesDB, error)
	Update(filter, update bson.M) (string, error)
	DeleteOne(filter bson.M) error
	Aggregate(pipeline bson.A) ([]entity.BookedBikesDB, error)
}
