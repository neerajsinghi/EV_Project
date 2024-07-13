package charger

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/charger/chargeDB"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var service = chargeDB.NewService()

// AddCharger adds a new station
func AddCharger(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	station, err := getCharger(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get station"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddCharger(station)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to add station"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to add station"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

// UpdateCharger updates a station
func UpdateCharger(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	station, err := getCharger(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get station"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateCharger(id, station)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to update station"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to update station"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

// DeleteCharger deletes a station
func DeleteCharger(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := mux.Vars(r)["id"]
	err := service.DeleteCharger(id)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to delete station"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": ""})
}

// GetAllChargers gets all stations
func GetAllChargers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := service.GetCharger()
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to get stations"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

// GetNearByChargers gets all stations near by
func GetNearByChargers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("long")
	distance := r.URL.Query().Get("distance")
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)
	distanceInt, _ := strconv.Atoi(distance)
	data, err := service.GetNearByCharger(latFloat, lonFloat, distanceInt)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to get stations"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

func getCharger(r *http.Request) (entity.ChargerDB, error) {
	var station entity.ChargerDB
	err := json.NewDecoder(r.Body).Decode(&station)
	if err != nil {
		return entity.ChargerDB{}, errors.Wrap(err, "unable to decode request payload")
	}
	return station, nil
}
