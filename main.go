package main

import (
	"log"

	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/DevMode"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Controller"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Model"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Router"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Utility"
)

func main() {

	Util, err := Utility.NewUtility(DevMode.Client, "./")
	if err != nil {
		log.Fatal(err)
	}

	Model, err := Model.NewModel(Util)
	if err != nil {
		Util.DB.Close()
		log.Fatal(err)
	}

	Ctrl, err := Controller.NewController(Util, Model)

	if err != nil {
		Util.DB.Close()
		log.Fatal(err)
	}

	router, err := Router.NewRouter(Ctrl)
	if err != nil {
		Util.DB.Close()
		log.Fatal(err)
	}

	if err := router.Run("0.0.0.0:8090"); err != nil {
		Util.DB.Close()
		log.Fatal(err)
	}
}
