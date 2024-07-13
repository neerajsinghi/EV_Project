package wallet

import (
	"bikeRental/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	InsertOne(document interface{}) (string, error)
	Find(filter, projection bson.M) ([]entity.WalletS, error)
	DeleteOne(filter bson.M) error
}
