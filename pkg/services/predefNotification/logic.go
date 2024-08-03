package predefnotification

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/notificationModel"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = notificationModel.NewRepository("predefNotification")

// InsertOne inserts a new predefNotification into the database
func InsertOne(predefNotification interface{}) (string, error) {
	res, err := repo.InsertOne(predefNotification)
	if err != nil {
		return "", errors.New("error in inserting predefNotification")
	}
	return res, nil
}

// Get gets a predefNotification from the database based on name
func Get(name string) (entity.PreDefNotification, error) {
	data, err := repo.FindOne(bson.M{"name": name}, nil)
	if err != nil {
		return entity.PreDefNotification{}, errors.New("error in finding predefNotification")
	}
	return data, nil
}

// GetAll gets all predefNotifications from the database
func GetAll() ([]entity.PreDefNotification, error) {
	data, err := repo.Find(bson.M{}, nil)
	if err != nil {
		return nil, errors.New("error in finding predefNotification")
	}
	return data, nil
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
	data, err := repo.UpdateOne(bson.M{"name": name}, bson.M{"$set": set})
	if err != nil {
		return "", errors.New("error in updating predefNotification")
	}
	return data, nil
}

// DeleteOne deletes a predefNotification from the database
func DeleteOne(name string) error {
	err := repo.DeleteOne(bson.M{"name": name})
	if err != nil {
		return errors.New("error in deleting predefNotification")
	}
	return nil
}
