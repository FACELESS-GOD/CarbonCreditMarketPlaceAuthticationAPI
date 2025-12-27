package Controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/DevMode"
	HeaderCaption "github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/Header"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Model"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Utility"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type TestControllerStruct struct {
	Ut   Utility.Utils
	Mdl  Model.ModelStruct
	Ctrl ControllerStruct
	suite.Suite
	VerEmail      string
	VerPass       string
	UserID        int
	RefereshToken string
	AccessToken   string
}

func (Ts *TestControllerStruct) Reset() {
	Ts.VerEmail = ""
	Ts.VerPass = ""
	Ts.UserID = 0
	Ts.RefereshToken = ""
	Ts.AccessToken = ""

}

func TestMain(m *testing.T) {
	suite.Run(m, &TestControllerStruct{})
}

func (Ts *TestControllerStruct) SetupSuite() {

	util, err := Utility.NewUtility(DevMode.Test, "../../")
	if err != nil {
		Ts.FailNow(err.Error())
	}
	Ts.Ut = util

	mdl, err := Model.NewModel(Ts.Ut)
	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Mdl = mdl

	ctrl, err := NewController(util, mdl)
	if err != nil {
		Ts.FailNow(err.Error())
	}
	Ts.Ctrl = ctrl
}

func (Its *TestControllerStruct) TestRefereshToken() {

	router := gin.Default()

	router.GET("/tkn", Its.Ctrl.AuthMiddleWare(), Its.Ctrl.RefereshToken)

	recorder := httptest.NewRecorder()

	cookie := http.Cookie{
		Name:     "__host-http-Login",
		Value:    Its.RefereshToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(48 * time.Hour),
		MaxAge:   86400 * 2,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	req1, err := http.NewRequest("GET", "/tkn", nil)

	if err != nil {
		Its.FailNow(err.Error())
	}

	req1.AddCookie(&cookie)

	req1.Header.Add(HeaderCaption.Authorization, "Bearer "+Its.AccessToken)

	router.ServeHTTP(recorder, req1)

	Its.Require().Equal(200, recorder.Code)

	tkn := recorder.Result().Header.Get(HeaderCaption.Authorization)

	Its.Require().Equal(len(tkn) >= 1, true)

	cookies := recorder.Result().Cookies()
	Its.Require().Equal(len(cookies) >= 1, true)

	isPresent := false

	for _, cookie := range cookies {
		if cookie.Name == "__host-http-Login" {
			isPresent = true
			Its.Require().Equal(len(cookie.Value) >= 1, true)

		}
	}
	Its.Require().Equal(isPresent, true)

}

func (Its *TestControllerStruct) TestVerifyCred() {

	router := gin.Default()

	router.GET("/login", Its.Ctrl.VerifyCred)

	recorder := httptest.NewRecorder()

	currCred := Cred{
		EMAIL:    Its.VerEmail,
		PASSWORD: Its.VerPass,
	}
	jsonData, err := json.Marshal(&currCred)

	if err != nil {
		Its.FailNow(err.Error())
	}

	req1, _ := http.NewRequest("GET", "/login", strings.NewReader(string(jsonData)))

	if err != nil {
		Its.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Its.Require().Equal(200, recorder.Code)

	tkn := recorder.Result().Header.Get(HeaderCaption.Authorization)

	Its.Require().Equal(len(tkn) >= 1, true)

	cookies := recorder.Result().Cookies()
	Its.Require().Equal(len(cookies) >= 1, true)

	isPresent := false

	for _, cookie := range cookies {
		if cookie.Name == "__host-http-Login" {
			isPresent = true
			Its.Require().Equal(len(cookie.Value) >= 1, true)

		}
	}
	Its.Require().Equal(isPresent, true)

}

func (Its *TestControllerStruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestVerifyCred":
		Its.Reset()
		randNum, err := Its.Ut.RandomNumber(10)
		if err != nil {
			Its.FailNow(err.Error())
		}
		base, err := Its.Ut.RandomString(int(randNum))

		if err != nil {
			Its.FailNow(err.Error())
		}

		tail, err := Its.Ut.RandomString(int(randNum))

		if err != nil {
			Its.FailNow(err.Error())
		}

		pass, err := Its.Ut.RandomString(int(randNum))

		if err != nil {
			Its.FailNow(err.Error())
		}

		Its.VerEmail = base + "@" + tail + ".com"
		Its.VerPass = pass

		err = Its.addUser()

		if err != nil {
			Its.FailNow(err.Error())
		}

	case "TestRefereshToken":
		Its.Reset()
		randNum, err := Its.Ut.RandomNumber(10)
		if err != nil {
			Its.FailNow(err.Error())
		}
		base, err := Its.Ut.RandomString(int(randNum))

		if err != nil {
			Its.FailNow(err.Error())
		}

		tail, err := Its.Ut.RandomString(int(randNum))

		if err != nil {
			Its.FailNow(err.Error())
		}

		pass, err := Its.Ut.RandomString(int(randNum))

		if err != nil {
			Its.FailNow(err.Error())
		}

		Its.VerEmail = base + "@" + tail + ".com"
		Its.VerPass = pass

		err = Its.addUser()

		if err != nil {
			Its.FailNow(err.Error())
		}

		isValid, refereshToken, accessToken, err := Its.Mdl.VerifyCred(Model.LoginStruct{Email: Its.VerEmail, Password: Its.VerPass})
		if err != nil {
			Its.FailNow(err.Error())
		}
		if isValid != true {
			Its.FailNow("Setup Failed")
		}
		if len(refereshToken) < 1 || len(accessToken) < 1 {
			Its.FailNow("Setup Failed")
		}
		Its.RefereshToken = refereshToken
		Its.AccessToken = accessToken

	}
}

func (Its *TestControllerStruct) AfterTest(SuiteName string, TestName string) {

	switch TestName {
	case "TestVerifyCred":
		err := Its.deleteUser()

		if err != nil {
			Its.FailNow(err.Error())
		}
	case "TestRefereshToken":

		err := Its.deleteUser()

		if err != nil {
			Its.FailNow(err.Error())
		}
	}
}

func (Its *TestControllerStruct) TearDownSuite() {
	Its.Ut.DB.Close()
}

func (Its *TestControllerStruct) GenerateHash(Password string) (string, error) {

	var customCost int = 15
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), customCost)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(hashedPassword), nil
}

const AddUserQuery string = `
INSERT INTO User (
  Name, email
) VALUES (
  ? , ? 
)
;
`

const AddUserCredQuery string = `
INSERT INTO UserCred (
  UserId, Hash_Password
) VALUES (
  ? , ? 
)
;
`

func (Its *TestControllerStruct) addUser() error {

	ctx := context.WithoutCancel(context.Background())

	db, err := Its.Ut.DB.BeginTx(ctx, &Its.Ut.TxOption)

	defer db.Commit()

	if err != nil {
		return err
	}

	name, err := Its.Ut.RandomString(int(10))

	if err != nil {
		Its.FailNow(err.Error())
	}

	response, err := db.ExecContext(ctx, AddUserQuery, name, Its.VerEmail)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	userID, err := response.LastInsertId()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	password, err := Its.GenerateHash(Its.VerPass)
	if err != nil {
		return err
	}

	response, err = db.ExecContext(ctx, AddUserCredQuery, userID, password)

	if err != nil {
		return err
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	Its.UserID = int(userID)
	return nil
}

const DeleteUserQuery string = `
UPDATE User
SET Is_Visible = 0 , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

const DeleteUserCredQuery string = `
UPDATE UserCred
SET Is_Visible = 0 , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Its *TestControllerStruct) deleteUser() error {

	ctx := context.WithoutCancel(context.Background())

	db, err := Its.Ut.DB.BeginTx(ctx, &Its.Ut.TxOption)
	defer db.Commit()
	if err != nil {
		return err
	}

	userCredresponse, err := db.ExecContext(ctx, DeleteUserCredQuery, Its.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	rowsaffected, err := userCredresponse.RowsAffected()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	if rowsaffected < 1 {
		return errors.New("User Doesnot Exsists")
	}

	userResponse, err := db.Query(DeleteUserQuery, Its.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	err = userResponse.Close()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			return errors.Join(nerr, err)
		} else {
			return err
		}
	}

	return nil

}
