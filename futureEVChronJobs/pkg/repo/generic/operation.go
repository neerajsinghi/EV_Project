package generic

import (
	"errors"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	CollectionName string
}

// NewFirestoreRepository creates a new repo
func NewRepository(collectionName string) BookingRepository {
	return &repo{
		CollectionName: collectionName,
	}
}

// used by signup
func (r *repo) InsertOne(document interface{}) (string, error) {
	user, err := trestCommon.InsertOne(document, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"insert data",
			err,
			logrus.Fields{
				"document":        document,
				"collection name": r.CollectionName,
			})
		return "", err
	}
	userid := user.InsertedID.(primitive.ObjectID).Hex()
	return userid, nil
}

// used by update data ,login and email verifcation
func (r *repo) UpdateOne(filter, update bson.M) (string, error) {
	result, err := trestCommon.UpdateOne(filter, update, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"update data",
			err,
			logrus.Fields{
				"filter":          filter,
				"update":          update,
				"collection name": r.CollectionName,
			})

		return "", err
	}
	if result.MatchedCount == 0 || result.ModifiedCount == 0 {
		err = errors.New("data not found(404)")
		trestCommon.ECLog3(
			"update data",
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

// used by get data ,login and email verification
func (r *repo) FindOne(filter, projection bson.M) (*mongo.Cursor, error) {
	return trestCommon.FindSort(filter, projection, bson.M{"created_time": -1}, 1, 0, r.CollectionName)
}

// not used may use in future for gettin list of datas
func (r *repo) Find(filter, projection bson.M) (*mongo.Cursor, error) {
	return trestCommon.FindSort(filter, projection, bson.M{"_id": -1}, 1000000, 0, r.CollectionName)
}
func (r *repo) Aggregate(pipeline bson.A) (*mongo.Cursor, error) {
	return trestCommon.Aggregate(pipeline, r.CollectionName)
}

// not using
func (r *repo) DeleteOne(filter bson.M) error {
	deleteResult, err := trestCommon.DeleteOne(filter, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"delete data",
			err,
			logrus.Fields{
				"filter":          filter,
				"collection name": r.CollectionName,
			})
		return err
	}
	if deleteResult.DeletedCount == 0 {
		err = errors.New("data not found(404)")
		trestCommon.ECLog3(
			"delete data",
			err,
			logrus.Fields{
				"filter":          filter,
				"collection name": r.CollectionName,
			})
		return err
	}
	return nil
}
