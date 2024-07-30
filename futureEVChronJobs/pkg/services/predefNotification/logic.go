package predefnotification

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/notificationModel"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = notificationModel.NewRepository("predefNotification")

// InsertOne inserts a new predefNotification into the database
func InsertOne(predefNotification interface{}) (string, error) {
	return repo.InsertOne(predefNotification)
}

// Get gets a predefNotification from the database based on name
func Get(name string) (entity.PreDefNotification, error) {
	return repo.FindOne(bson.M{"name": name}, nil)
}

// GetAll gets all predefNotifications from the database
func GetAll() ([]entity.PreDefNotification, error) {
	return repo.Find(bson.M{}, nil)
}

// UpdateOne updates a predefNotification in the database
func UpdateOne(name string, update entity.PreDefNotification) (string, error) {
	set := bson.M{}

	if update.Title != "" {
		set["title"] = update.Title
	}
	if update.Body != "" {
		set["body"] = update.Body
	}
	if update.Type != "" {
		set["type"] = update.Type
	}
	return repo.UpdateOne(bson.M{"name": name}, bson.M{"$set": set})
}

// DeleteOne deletes a predefNotification from the database
func DeleteOne(name string) error {
	return repo.DeleteOne(bson.M{"name": name})
}
