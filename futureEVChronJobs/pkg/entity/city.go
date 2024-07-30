package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// ProfileDB ...
type City struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Name             string             `json:"name" bson:"name"`
	Active           *bool              `json:"active" bson:"active"`
	NumberOfStations *int               `json:"numberOfStations" bson:"numberOfStations"`
	NumberOfVehicles *int               `json:"numberOfVehicles" bson:"numberOfVehicles"`
	LocationPolygon  LocationPolygon    `json:"locationPolygon" bson:"locationPolygon"`
	Services         []string           `json:"services" bson:"services"`
	VehicleType      string             `json:"vehicleType" bson:"vehicleType"`
}
type LocationPolygon struct {
	Type        string          `json:"type" bson:"type"`
	Coordinates [][][][]float64 `json:"coordinates" bson:"coordinates"`
}
