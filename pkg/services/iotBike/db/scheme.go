package db

import "bikeRental/pkg/entity"

type IOtBike interface {
	FindAll() ([]entity.IotBikeDB, error)
	FindNearByBikes(lat, long float64, distance int, bType string) ([]entity.IotBikeDB, error)
}
