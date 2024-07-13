package iotbike

import (
	"bikeRental/pkg/entity"
	"context"
	"errors"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repo struct {
	CollectionName string
}

// DeleteOne implements IOTBikeRepository.
func (r *repo) DeleteOne(filter primitive.M) error {
	panic("unimplemented")
}

// Find implements IOTBikeRepository.
func (r *repo) Find(filter primitive.M, projection primitive.M) ([]entity.IotBikeDB, error) {
	var profiles []entity.IotBikeDB
	cursor, err := trestCommon.FindSort(filter, projection, bson.M{"_id": -1}, 1000000, 0, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"Find profiles",
			err,
			logrus.Fields{
				"filter":          filter,
				"collection name": r.CollectionName,
			})
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.TODO()) {
		var profile entity.IotBikeDB
		if err = cursor.Decode(&profile); err != nil {
			trestCommon.ECLog3(
				"Find profiles",
				err,
				logrus.Fields{
					"filter":          filter,
					"collection name": r.CollectionName,
					"error at":        cursor.RemainingBatchLength(),
				})
			return profiles, nil
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil

}

// FindOne implements IOTBikeRepository.
func (r *repo) FindOne(filter primitive.M, projection primitive.M) (entity.IotBikeDB, error) {
	panic("unimplemented")
}

// InsertOne implements IOTBikeRepository.
func (r *repo) InsertOne(document interface{}) (string, error) {
	panic("unimplemented")
}

// UpdateOne implements IOTBikeRepository.
func (r *repo) UpdateOne(filter primitive.M, update primitive.M) (string, error) {
	result, err := trestCommon.UpdateOne(filter, update, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"update profile",
			err,
			logrus.Fields{
				"filter":          filter,
				"update":          update,
				"collection name": r.CollectionName,
			})

		return "", err
	}
	if (result.MatchedCount == 0 || result.ModifiedCount == 0) && result.UpsertedCount == 0 {
		err = errors.New("profile not found(404)")
		trestCommon.ECLog3(
			"update profile",
			err,
			logrus.Fields{
				"filter":          filter,
				"update":          update,
				"collection name": r.CollectionName,
			})
		return "", err
	}
	return "updated successfully", nil
}

// NewFirestoreRepository creates a new repo
func NewProfileRepository(collectionName string) IOTBikeRepository {
	return &repo{
		CollectionName: collectionName,
	}
}
