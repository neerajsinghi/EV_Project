package db

import (
	"futureEVChronJobs/pkg/entity"
	iotbike "futureEVChronJobs/pkg/repo/iot_bike"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	repo = iotbike.NewProfileRepository("iotBike")
)

func GetBike(ids []int) ([]entity.IotBikeDB, error) {
	filter := bson.A{}
	for _, id := range ids {
		filter = append(filter, bson.M{"deviceId": id})
	}
	return repo.Find(bson.M{"$or": filter}, bson.M{})
}
