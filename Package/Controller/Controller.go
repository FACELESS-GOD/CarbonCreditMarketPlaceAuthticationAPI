package Controller

import (
	"log"
	"net/http"
	"strings"
	"time"

	Contextmapcaption "github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/ContextMapCaption"
	HeaderCaption "github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/Header"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Model"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Utility"
	"github.com/gin-gonic/gin"
)

type ControllerStruct struct {
	Mdl Model.ModelStruct
	Ut  Utility.Utils
}

func NewController(Ut Utility.Utils, Model Model.ModelStruct) (ControllerStruct, error) {
	mdl := ControllerStruct{}
	mdl.Ut = Ut
	mdl.Mdl = Model
	return mdl, nil
}

type Cred struct {
	EMAIL    string `json:"EMAIL"`
	PASSWORD string `json:"PASSWORD"`
}

func (Ctrl *ControllerStruct) VerifyCred(gCtx *gin.Context) {
	currCred := Cred{}
	err := gCtx.BindJSON(&currCred)

	if err != nil {
		gCtx.JSON(http.StatusBadRequest, currCred)
		return
	}

	cred := Model.LoginStruct{
		Email:    currCred.EMAIL,
		Password: currCred.PASSWORD,
	}

	isValid, refereshToken, accessToken, err := Ctrl.Mdl.VerifyCred(cred)

	if err != nil {
		gCtx.Status(http.StatusInternalServerError)
		return
	}

	if isValid != true {
		gCtx.Status(http.StatusBadRequest)
		return
	}

	if len(accessToken) < 1 || len(refereshToken) < 1 {
		gCtx.Status(http.StatusInternalServerError)
		return
	}

	tokenString := "Bearer " + accessToken

	cookie := http.Cookie{
		Name:     "__host-http-Login",
		Value:    refereshToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(48 * time.Hour),
		MaxAge:   86400 * 2,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	gCtx.Header(HeaderCaption.Authorization, tokenString)
	gCtx.SetCookieData(&cookie)
	return

}

func (Ctrl *ControllerStruct) RefereshToken(gCtx *gin.Context) {
	userId, ok := gCtx.Keys[Contextmapcaption.UserId]

	if ok != true {
		gCtx.Status(http.StatusInternalServerError)
		return
	}

	Id, ok := userId.(int)
	if ok != true {
		gCtx.Status(http.StatusInternalServerError)
		return
	}

	refereshToken, accessToken, err := Ctrl.Mdl.AddRefereshTokenToDB(int(Id))

	if err != nil {
		gCtx.Status(http.StatusInternalServerError)
		return
	}

	if len(accessToken) < 1 || len(refereshToken) < 1 {
		gCtx.Status(http.StatusInternalServerError)
		return
	}

	tokenString := "Bearer " + accessToken

	cookie := http.Cookie{
		Name:     "__host-http-Login",
		Value:    refereshToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(48 * time.Hour),
		MaxAge:   86400 * 2,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	gCtx.Header(HeaderCaption.Authorization, tokenString)
	gCtx.SetCookieData(&cookie)
	return
}

func (Ctrl *ControllerStruct) CustomRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if err := recover(); err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func (Ctrl *ControllerStruct) AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader(HeaderCaption.Authorization)
		if len(tokenString) < 1 {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		tokenSlice := strings.Split(tokenString, " ")

		IsValid, tknm, err := Ctrl.Mdl.Ut.VerifyToken(tokenSlice[1])

		if IsValid != true {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid or expired token"})
			return
		}

		// if tknm.ExpiresAt.Compare(time.Now()) >= 1 {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		// 	return
		// }

		ctx.Set(Contextmapcaption.UserId, tknm.UserID)

		ctx.Next()
	}
}
