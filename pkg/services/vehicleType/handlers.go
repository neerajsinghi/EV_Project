package vehicletype

import (
	"bikeRental/pkg/entity"
	vdb "bikeRental/pkg/services/vehicleType/vDB"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var service = vdb.NewService()

// AddVehicleType adds a new vehicle type
func AddVehicleType(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	vehicleType, err := getVehicleType(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get vehicle type"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddVehicleType(vehicleType)
	utils.SendOutput(err, w, r, data, vehicleType, "AddVehicleType")
}

// UpdateVehicleType updates a vehicle type
func UpdateVehicleType(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	vehicleType, err := getVehicleType(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get vehicle type"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateVehicleType(id, vehicleType)
	utils.SendOutput(err, w, r, data, vehicleType, "UpdateVehicleType")
}

// DeleteVehicleType deletes a vehicle type
func DeleteVehicleType(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	err := service.DeleteVehicleType(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteVehicleType")
}

func GetAllVehicleTypes(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	data, err := service.GetVehicleType()
	utils.SendOutput(err, w, r, data, nil, "GetAllVehicleTypes")
}

func getVehicleType(r *http.Request) (entity.VehicleTypeDB, error) {
	var vehicleType entity.VehicleTypeDB
	err := json.NewDecoder(r.Body).Decode(&vehicleType)
	if err != nil {
		return entity.VehicleTypeDB{}, errors.Wrapf(err, "unable to decode request body")
	}
	return vehicleType, nil
}
