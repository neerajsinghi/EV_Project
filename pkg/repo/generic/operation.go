package generic

import (
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
func (r *repo) FindOne(filter, projection bson.M) (interface{}, error) {
	var data interface{}
	cursor, err := trestCommon.FindSort(filter, projection, bson.M{"created_time": -1}, 1, 0, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"Find data",
			err,
			logrus.Fields{
				"filter":          filter,
				"collection name": r.CollectionName,
			})
		return data, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(&data); err != nil {
			trestCommon.ECLog3(
				"Find datas",
				err,
				logrus.Fields{
					"filter":          filter,
					"collection name": r.CollectionName,
					"error at":        cursor.RemainingBatchLength(),
				})
			return data, nil
		}
		break
	}
	return data, err
}

// not used may use in future for gettin list of datas
func (r *repo) Find(filter, projection bson.M) ([]interface{}, error) {
	var datas []interface{}
	cursor, err := trestCommon.FindSort(filter, projection, bson.M{"_id": -1}, 1000000, 0, r.CollectionName)
	if err != nil {
		trestCommon.ECLog3(
			"Find datas",
			err,
			logrus.Fields{
				"filter":          filter,
				"collection name": r.CollectionName,
			})
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.TODO()) {
		var data interface{}
		if err = cursor.Decode(&data); err != nil {
			trestCommon.ECLog3(
				"Find datas",
				err,
				logrus.Fields{
					"filter":          filter,
					"collection name": r.CollectionName,
					"error at":        cursor.RemainingBatchLength(),
				})
			return datas, nil
		}
		datas = append(datas, data)
	}
	return datas, nil
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
