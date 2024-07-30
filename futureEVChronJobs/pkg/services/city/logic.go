package city

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/city"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = city.NewRepository("city")

func GetCity(id string) (*entity.City, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	return repo.FindOne(filter, bson.M{})
}
func GetCities() ([]entity.City, error) {
	return repo.Find(bson.M{}, bson.M{})
}

func InCity(lat, long float64) (*entity.City, error) {
	filter := bson.M{
		"locationPolygon": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}
	city, err := repo.FindOne(filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return city, nil
}
