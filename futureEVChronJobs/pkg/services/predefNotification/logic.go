package predefnotification

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/notificationModel"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = notificationModel.NewRepository("predefNotification")

// Get gets a predefNotification from the database based on name
func Get(name string) (entity.PreDefNotification, error) {
	return repo.FindOne(bson.M{"name": name}, nil)
}
