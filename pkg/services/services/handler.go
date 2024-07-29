package services

import (
	"bikeRental/pkg/entity"
	servdb "bikeRental/pkg/services/services/servDB"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var serviceDB = servdb.NewService()

// AddService adds a new service
func AddService(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	service, err := getService(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get service"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := serviceDB.InsertOne(service)
	utils.SendOutput(err, w, r, data, service, "AddService")
}

// UpdateService updates a service
func UpdateService(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	service, err := getService(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get service"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := serviceDB.UpdateService(id, service)
	utils.SendOutput(err, w, r, data, service, "UpdateService")
}

// DeleteService deletes a service
func DeleteService(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	err := serviceDB.DeleteService(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteService")
}

func GetService(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := serviceDB.GetAllServices()
	utils.SendOutput(err, w, r, data, nil, "GetAllServices")
}
func getService(r *http.Request) (entity.ServiceDB, error) {
	var service entity.ServiceDB
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		return entity.ServiceDB{}, errors.Wrapf(err, "unable to decode request body")
	}
	return service, nil
}
