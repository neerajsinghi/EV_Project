package vdb

import "bikeRental/pkg/entity"

type Serv interface {
	AddVehicleType(document entity.VehicleTypeDB) (string, error)
	UpdateVehicleType(id string, document entity.VehicleTypeDB) (string, error)
	DeleteVehicleType(id string) error
	GetVehicleType() ([]entity.VehicleTypeDB, error)
	GetVehicleTypeByID(id string) (entity.VehicleTypeDB, error)
}
