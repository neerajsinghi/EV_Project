package selectedcolumns

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/generic"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = generic.NewRepository("selectedColumns")

func SelectColumns(document entity.ColumnEntity) (*entity.ColumnEntity, error) {
	bsonDoc := bson.M{
		"user_id":          document.UserID,
		"table_name":       document.TableName,
		"columns_selected": document.ColumnsSelected,
	}
	repo.UpdateOne(bson.M{"user_id": document.UserID, "table_name": document.TableName}, bson.M{"$set": bsonDoc})

	return GetColumns(document.UserID, document.TableName)
}

func GetColumns(userId, tableName string) (*entity.ColumnEntity, error) {
	data, err := repo.FindOne(bson.M{"user_id": userId, "table_name": tableName}, nil)
	if err != nil {
		return nil, errors.New("error in finding columns")
	}
	defer data.Close(context.Background())
	var result []entity.ColumnEntity
	err = data.All(context.Background(), &result)
	if err != nil {
		return nil, errors.New("error in finding columns")
	}
	return &result[0], nil
}

func GetColumnsForUser(userId string) ([]entity.ColumnEntity, error) {
	data, err := repo.Find(bson.M{"user_id": userId}, nil)
	if err != nil {
		return nil, errors.New("error in finding columns")
	}
	defer data.Close(context.Background())
	var result []entity.ColumnEntity
	err = data.All(context.Background(), &result)
	if err != nil {
		return nil, errors.New("error in finding columns")
	}
	return result, nil
}

func DeleteColumns(userId, tableName string) error {
	err := repo.DeleteOne(bson.M{"user_id": userId, "table": tableName})
	if err != nil {
		return errors.New("error in deleting columns")
	}
	return err
}
