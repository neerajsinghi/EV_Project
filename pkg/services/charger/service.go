package charger

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/charger/chargeDB"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"
	"strconv"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var service = chargeDB.NewService()

// AddCharger adds a new station
func AddCharger(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	station, err := getCharger(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get station"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddCharger(station)
	utils.SendOutput(err, w, r, data, station, "AddCharger")
}

// UpdateCharger updates a station
func UpdateCharger(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	station, err := getCharger(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get station"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateCharger(id, station)
	utils.SendOutput(err, w, r, data, station, "UpdateCharger")
}

// DeleteCharger deletes a station
func DeleteCharger(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	id := mux.Vars(r)["id"]
	err := service.DeleteCharger(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteCharger")
}

// GetAllChargers gets all stations
func GetAllChargers(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	data, err := service.GetCharger()
	utils.SendOutput(err, w, r, data, nil, "GetAllChargers")
}

// GetNearByChargers gets all stations near by
func GetNearByChargers(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("long")
	distance := r.URL.Query().Get("distance")
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)
	distanceInt, _ := strconv.Atoi(distance)
	data, err := service.GetNearByCharger(latFloat, lonFloat, distanceInt)
	utils.SendOutput(err, w, r, data, nil, "GetNearByChargers")
}

func getCharger(r *http.Request) (entity.ChargerDB, error) {
	var station entity.ChargerDB
	err := json.NewDecoder(r.Body).Decode(&station)
	if err != nil {
		return entity.ChargerDB{}, errors.Wrap(err, "unable to decode request payload")
	}
	return station, nil
}
