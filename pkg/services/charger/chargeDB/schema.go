package chargeDB

import "bikeRental/pkg/entity"

type Serv interface {
	AddCharger(document entity.ChargerDB) (string, error)
	UpdateCharger(id string, document entity.ChargerDB) (string, error)
	DeleteCharger(id string) error
	GetCharger() ([]entity.ChargerDB, error)
	GetChargerByID(id string) (entity.ChargerDB, error)
	GetNearByCharger(lat, long float64, distance int) ([]entity.ChargerDB, error)
}
