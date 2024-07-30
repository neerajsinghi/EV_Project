package db

import (
	"futureEVChronJobs/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type Wallet interface {
	InsertOne(document entity.WalletS) (WalletTotal, error)
	FindMy(userId string) (WalletTotal, error)
	Find() ([]WalletTotal, error)
	DeleteOne(filter bson.M) error
}

type WalletTotal struct {
	Wallets         []entity.WalletS
	TotalBalance    float64
	RefundableMoney float64
	UserData        *entity.ProfileDB
}
