package city

import (
	"encoding/json"
	"futureEVChronJobs/pkg/entity"
	utils "futureEVChronJobs/pkg/util"
	"net/http"
	"strconv"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func AddCityHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	booking, err := parseCity(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := AddCity(booking)
	utils.SendOutput(err, w, r, data, nil, "AddCity")
}

func GetAllCitiesHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := GetCities()
	utils.SendOutput(err, w, r, data, nil, "GetAllCities")
}
func GetCityHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	data, err := GetCity(id)
	utils.SendOutput(err, w, r, data, nil, "GetCity")
}
func InCityHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	lat := r.URL.Query().Get("lat")
	long := r.URL.Query().Get("long")
	latF, _ := strconv.ParseFloat(lat, 64)
	longF, _ := strconv.ParseFloat(long, 64)
	data, err := InCity(latF, longF)
	utils.SendOutput(err, w, r, data, nil, "InCity")
}
func UpdateCityHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	city, err := parseCity(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := UpdateCity(id, city)
	utils.SendOutput(err, w, r, data, city, "UpdateCity")
}

func DeleteCityHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	err := DeleteCity(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteCity")
}

func parseCity(r *http.Request) (entity.City, error) {
	var city entity.City
	if r.Body == nil {
		return city, nil
	}
	err := json.NewDecoder(r.Body).Decode(&city)
	return city, err
}
