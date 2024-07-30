package station

import (
	"bikeRental/pkg/entity"
	sdb "bikeRental/pkg/services/station/sDB"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"
	"strconv"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var service = sdb.NewService()

// AddStation adds a new station
func AddStation(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	station, err := getStation(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get station"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddStation(station)
	utils.SendOutput(err, w, r, data, station, "AddStation")
}

// UpdateStation updates a station
func UpdateStation(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	station, err := getStation(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get station"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateStation(id, station)
	utils.SendOutput(err, w, r, data, station, "UpdateStation")
}

// DeleteStation deletes a station
func DeleteStation(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	id := mux.Vars(r)["id"]
	err := service.DeleteStation(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteStation")
}

// GetAllStations gets all stations
func GetAllStations(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	data, err := service.GetStation()
	utils.SendOutput(err, w, r, data, nil, "GetAllStations")
}

// GetNearByStations gets all stations near by
func GetNearByStations(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("long")
	distance := r.URL.Query().Get("distance")
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)
	distanceInt, _ := strconv.Atoi(distance)
	data, err := service.GetNearByStation(latFloat, lonFloat, distanceInt)
	utils.SendOutput(err, w, r, data, nil, "GetNearByStations")
}
func GetStationsByID(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	id := mux.Vars(r)["id"]
	data, err := service.GetStationByID(id)
	utils.SendOutput(err, w, r, data, nil, "GetStationsByID")
}

func getStation(r *http.Request) (entity.StationDB, error) {
	var station entity.StationDB
	err := json.NewDecoder(r.Body).Decode(&station)
	if err != nil {
		return entity.StationDB{}, errors.Wrap(err, "unable to decode request payload")
	}
	return station, nil
}
