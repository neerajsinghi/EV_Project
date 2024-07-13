package sdb

import "bikeRental/pkg/entity"

type Serv interface {
	AddStation(document entity.StationDB) (string, error)
	UpdateStation(id string, document entity.StationDB) (string, error)
	DeleteStation(id string) error
	GetStation() ([]entity.StationDB, error)
	GetStationByID(id string) (entity.StationDB, error)
	GetNearByStation(lat, long float64, distance int) ([]entity.StationDB, error)
}
