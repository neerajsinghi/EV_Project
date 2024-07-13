package db

import "bikeRental/pkg/entity"

type IOtBike interface {
	AddBikeDevice(document entity.DeviceInfo) (string, error)
	UpdateBikeDevice(id string, document entity.DeviceInfo) (string, error)
	DeleteBikeDevice(id string) error
	FindAll() ([]entity.DeviceInfo, error)
	FindBikeByStation(stationID string) ([]entity.DeviceInfo, error)
	FindBikeByDeviceID(deviceId string) ([]entity.DeviceInfo, error)
}
