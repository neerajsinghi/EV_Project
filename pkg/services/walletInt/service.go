package wallet

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/wallet/db"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var service = db.NewService()

func AddWallet(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	wallet, err := getWallet(r)
	utils.SendOutput(err, w, r, wallet, "AddWallet")
}
func GetAllWallets(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := service.Find()
	utils.SendOutput(err, w, r, data, "GetAllWallets")
}
func GetMyWallet(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	data, err := service.FindMy(id)
	utils.SendOutput(err, w, r, data, "GetMyWallet")
}
func getWallet(r *http.Request) (entity.WalletS, error) {
	var wallet entity.WalletS
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}
