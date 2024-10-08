package wallet

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/wallet/db"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"
	"strings"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
)

var service = db.NewService()

func AddWallet(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	wallet, err := getWallet(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.InsertOne(wallet)
	utils.SendOutput(err, w, r, data, wallet, "AddWallet")
}
func GetAllWallets(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := service.Find()
	utils.SendOutput(err, w, r, data, nil, "GetAllWallets")
}
func GetMyWallet(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	data, err := service.FindMy(id)
	utils.SendOutput(err, w, r, data, nil, "GetMyWallet")
}

func GetWalletByPlanID(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	if strings.EqualFold(id, "all") {
		id = ""
	}
	data, err := service.FindForPlan(id)
	utils.SendOutput(err, w, r, data, nil, "GetWalletByPlayerID")
}
func getWallet(r *http.Request) (entity.WalletS, error) {
	var wallet entity.WalletS
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}
