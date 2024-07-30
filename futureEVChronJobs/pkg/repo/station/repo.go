package station

import (
	"futureEVChronJobs/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	InsertOne(document interface{}) (string, error)
	FindOne(filter, projection bson.M) (entity.StationDB, error)
	Find(filter, projection bson.M) ([]entity.StationDB, error)
	UpdateOne(filter, update bson.M) (string, error)
	DeleteOne(filter bson.M) error
}
