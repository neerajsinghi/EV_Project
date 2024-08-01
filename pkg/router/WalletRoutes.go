package router

import wallet "bikeRental/pkg/services/wallet"

var walletRoutes = Routes{
	Route{
		"Add Wallet",
		"POST",
		"/wallet",
		wallet.AddWallet,
	},
	Route{
		"Get Wallet",
		"GET",
		"/wallet/{id}",
		wallet.GetMyWallet,
	},
	Route{
		"Get Wallet All",
		"GET",
		"/wallet",
		wallet.GetAllWallets,
	},
	Route{
		"Get Wallet By Plan",
		"GET",
		"/wallet/plan/{id}",
		wallet.GetWalletByPlanID,
	},
}
