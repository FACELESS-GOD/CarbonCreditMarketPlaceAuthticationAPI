package router

import (
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(Ctrl Controller.ControllerStruct) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(Ctrl.CustomRecovery())
	router.GET("/login", Ctrl.VerifyCred)
	router.GET("/tkn", Ctrl.AuthMiddleWare(), Ctrl.RefereshToken)
	return router, nil
}
