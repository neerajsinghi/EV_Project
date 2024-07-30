package notificationModel

import (
	"futureEVChronJobs/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type PlanRepository interface {
	InsertOne(document interface{}) (string, error)
	FindOne(filter, projection bson.M) (entity.PreDefNotification, error)
	Find(filter, projection bson.M) ([]entity.PreDefNotification, error)
	UpdateOne(filter, update bson.M) (string, error)
	DeleteOne(filter bson.M) error
}
