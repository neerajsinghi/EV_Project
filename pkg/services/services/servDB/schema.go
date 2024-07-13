package servdb

import "bikeRental/pkg/entity"

type IService interface {
	InsertOne(document entity.ServiceDB) (string, error)
	GetAllServices() ([]entity.ServiceDB, error)
	UpdateService(id string, document entity.ServiceDB) (string, error)
	DeleteService(id string) error
}
