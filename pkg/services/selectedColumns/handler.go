package selectedcolumns

import (
	"bikeRental/pkg/entity"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func SelectColumnsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	columns, err := getColumns(r)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := SelectColumns(columns)
	utils.SendOutput(err, w, r, data, columns, "SelectColumns")
}

func GetAllSelectedColumnsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	userID := mux.Vars(r)["id"]
	tableName := mux.Vars(r)["table"]
	data, err := GetColumns(userID, tableName)

	utils.SendOutput(err, w, r, data, nil, "GetAllSelectedColumns")
}
func GetSelectedColumnsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	userID := mux.Vars(r)["id"]
	data, err := GetColumnsForUser(userID)
	utils.SendOutput(err, w, r, data, nil, "GetSelectedColumns")
}
func DeleteColumnsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	userID := mux.Vars(r)["id"]
	tableName := mux.Vars(r)["table"]
	err := DeleteColumns(userID, tableName)
	utils.SendOutput(err, w, r, nil, nil, "DeleteColumns")
}
func getColumns(r *http.Request) (entity.ColumnEntity, error) {
	var columns entity.ColumnEntity
	err := json.NewDecoder(r.Body).Decode(&columns)
	if err != nil {
		return columns, err
	}
	return columns, nil
}
