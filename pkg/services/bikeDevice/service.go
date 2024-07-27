package bikeDevice

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/bikeDevice/db"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var service = db.NewService()

// AddBikeDevice adds a new bike
func AddBikeDevice(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	bikeDevice, err := getBikeDevice(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get bike device"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddBikeDevice(bikeDevice)
	utils.SendOutput(err, w, r, data, "AddBikeDevice")
}

// UpdateBikeDevice updates a bike
func UpdateBikeDevice(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	bikeDevice, err := getBikeDevice(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get bike device"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateBikeDevice(id, bikeDevice)
	utils.SendOutput(err, w, r, data, "UpdateBikeDevice")
}

// DeleteBikeDevice deletes a bike
func DeleteBikeDevice(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	err := service.DeleteBikeDevice(id)
	utils.SendOutput(err, w, r, "Deleted successfully", "DeleteBikeDevice")
}

// FindAll returns all the bikes
func GetAll(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := service.FindAll()
	utils.SendOutput(err, w, r, data, "GetAll")

}
func GetBikeDevicesByStation(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	stationID := mux.Vars(r)["stationID"]
	data, err := service.FindBikeByStation(stationID)
	utils.SendOutput(err, w, r, data, "GetBikeDevicesByStation")
}
func GetBikeDevicesByDeviceID(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	stationID := mux.Vars(r)["id"]
	data, err := service.FindBikeByDeviceID(stationID)
	utils.SendOutput(err, w, r, data, "GetBikeDevicesByDeviceID")
}
func getBikeDevice(r *http.Request) (bikeDevice entity.DeviceInfo, err error) {
	err = json.NewDecoder(r.Body).Decode(&bikeDevice)
	if err != nil {
		return bikeDevice, errors.Wrapf(err, "unable to decode request body")
	}
	return bikeDevice, nil
}
