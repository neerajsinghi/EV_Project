package plan

import (
	"bikeRental/pkg/entity"
	pdb "bikeRental/pkg/services/plan/pDB"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var service = pdb.NewService()

func AddPlan(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	plan, err := getPlan(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddPlan(plan)
	utils.SendOutput(err, w, r, data, plan, "AddPlan")
}
func GetAllPlans(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	pType := r.URL.Query().Get("type")
	city := r.URL.Query().Get("city")
	data, err := service.GetPlans(pType, city)
	utils.SendOutput(err, w, r, data, nil, "GetAllPlans")
}
func GetPlanById(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	params := mux.Vars(r)
	data, err := service.GetPlan(params["id"])
	utils.SendOutput(err, w, r, data, nil, "GetPlanById")
}
func UpdatePlan(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	plan, err := getPlan(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.UpdatePlan(id, plan)
	utils.SendOutput(err, w, r, data, plan, "UpdatePlan")
}
func DeletePlan(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	err := service.DeletePlan(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeletePlan")
}
func getPlan(r *http.Request) (entity.PlanDB, error) {
	var plan entity.PlanDB
	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		return entity.PlanDB{}, err
	}
	return plan, nil
}

func GetDeposit(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	pType := r.URL.Query().Get("type")
	city := r.URL.Query().Get("city")
	data, err := pdb.GetDeposit(city, pType)
	utils.SendOutput(err, w, r, data, nil, "GetDeposit")
}
