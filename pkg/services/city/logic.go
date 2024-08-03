package city

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/city"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = city.NewRepository("city")

func AddCity(city entity.City) (string, error) {
	city.ID = primitive.NewObjectID()
	status, err := repo.InsertOne(city)
	if err != nil {
		return "", errors.New("error in inserting city")
	}
	return status, nil
}

func DeleteCity(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	err := repo.DeleteOne(bson.M{"_id": idObject})
	if err != nil {
		return errors.New("error in deleting city")
	}
	return nil
}

func GetCity(id string) (*entity.City, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	data, err := repo.FindOne(filter, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding cities")
	}
	return data, nil
}
func GetCities() ([]entity.City, error) {
	data, err := repo.Find(bson.M{}, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding cities")
	}
	return data, nil
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
		return nil, errors.New("error in finding city")
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
	data, err := repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
	if err != nil {
		return "", errors.New("error in updating city")
	}
	return data, nil
}
