package city

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/city"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = city.NewRepository("city")

func AddCity(city entity.City) (string, error) {
	city.ID = primitive.NewObjectID()
	return repo.InsertOne(city)
}

func DeleteCity(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

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

func UpdateCity(id string, city entity.City) (string, error) {
	set := bson.M{}
	if city.Active != nil {
		set["active"] = city.Active
	}
	if city.Name != "" {
		set["name"] = city.Name
	}
	if city.NumberOfStations != nil {
		set["numberOfStations"] = city.NumberOfStations
	}
	if city.NumberOfVehicles != nil {
		set["numberOfVehicles"] = city.NumberOfVehicles
	}
	if city.LocationPolygon.Type != "" {
		set["locationPolygon.type"] = city.LocationPolygon.Type
	}
	if len(city.LocationPolygon.Coordinates) != 0 {
		set["locationPolygon.coordinates"] = city.LocationPolygon.Coordinates
	}
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
}
