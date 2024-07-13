package db

import (
	"bikeRental/pkg/entity"
	iotbike "bikeRental/pkg/repo/iot_bike"

	"go.mongodb.org/mongo-driver/bson"
)

type service struct{}

func (s *service) FindAll() ([]entity.IotBikeDB, error) {
	filter := bson.M{}
	projection := bson.M{}
	return repo.Find(filter, projection)

}
func (s *service) FindNearByBikes(lat, long float64, distance int, bType string) ([]entity.IotBikeDB, error) {
	filter := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": distance,
			},
		},
	}
	if bType != "" && bType != "all" {
		filter["type"] = bType
	}
	projection := bson.M{}
	return repo.Find(filter, projection)
}

var (
	repo = iotbike.NewProfileRepository("iotBike")
)

// NewService creates a new service
func NewService() IOtBike {
	return &service{}
}

func GetBike(ids []int) ([]entity.IotBikeDB, error) {
	filter := bson.A{}
	for _, id := range ids {
		filter = append(filter, bson.M{"deviceId": id})
	}
	return repo.Find(bson.M{"$or": filter}, bson.M{})
}
