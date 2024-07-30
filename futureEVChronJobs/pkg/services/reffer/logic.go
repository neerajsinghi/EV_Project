package reffer

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/reffer"

	"go.mongodb.org/mongo-driver/bson"
)

var repo = reffer.NewRepository("referral")

func GetReffer() ([]entity.ReferralDB, error) {
	return repo.Find(bson.M{}, bson.M{})
}
