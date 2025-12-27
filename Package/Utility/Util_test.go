package Utility

import (
	"fmt"
	"testing"

	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/DevMode"
	"github.com/stretchr/testify/suite"
)

type TestTokenStruct struct {
	Ut Utils
	suite.Suite
	UserId int
	Token  string
}

func (Ts *TestTokenStruct) Reset() {
	Ts.UserId = 0
	Ts.Token = ""
}

func TestTokenstruct(testObj *testing.T) {
	suite.Run(testObj, &TestTokenStruct{})
}

func (Its *TestTokenStruct) SetupSuite() {
	util, err := NewUtility(DevMode.Test, "../../")
	if err != nil {
		Its.FailNow(err.Error())
	}
	Its.Ut = util

}

func (Its *TestTokenStruct) TestVerifyToken() {
	_, _, err := Its.Ut.VerifyToken("")
	Its.Require().NotNil(err)

	var nilId string

	_, _, err = Its.Ut.VerifyToken(nilId)
	Its.Require().NotNil(err)

	_, _, err = Its.Ut.VerifyToken("iawsujdgfbh")
	Its.Require().NotNil(err)

	isValid, TokenMetadt, err := Its.Ut.VerifyToken(Its.Token)
	Its.Require().Nil(err)
	Its.Require().NotNil(TokenMetadt)
	Its.Require().Equal(isValid, true)
	Its.Require().Equal(len(TokenMetadt.TokenID.String()) >= 1, true)
	Its.Require().Equal(TokenMetadt.UserID >= 1, true)
	Its.Require().Equal(TokenMetadt.ExpiredAt.After(TokenMetadt.IssuedAT), true)

}
func (Its *TestTokenStruct) TestCreateToken() {
	_, err := Its.Ut.CreateToken(0, "Referesh")
	Its.Require().NotNil(err)

	var nilId int

	_, err = Its.Ut.CreateToken(nilId, "Referesh")
	Its.Require().NotNil(err)

	tkn, err := Its.Ut.CreateToken(Its.UserId, "Referesh")
	Its.Require().Nil(err)
	Its.Require().Equal(len(tkn) >= 1, true)

}

func (Its *TestTokenStruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestCreateToken":
		Its.Reset()
		Its.UserId = 10
	case "TestVerifyToken":
		Its.Reset()
		Its.UserId = 10
		tkn, err := Its.Ut.CreateToken(Its.UserId, "Referesh")
		if err != nil {
			Its.FailNow(err.Error())
		}

		Its.Token = tkn

	}
}

func (Its *TestTokenStruct) AfterTest(SuiteName string, TestName string) {

	switch TestName {
	case "TestCreateToken":
		Its.Reset()
	case "TestVerifyToken":
		Its.Reset()
	}
}

func (Its *TestTokenStruct) TearDownSuite() {
	Its.Ut.DB.Close()
}

type TestUtil struct {
	suite.Suite
	ut Utils
}

func TestUtils(t *testing.T) {
	suite.Run(t, &TestUtil{})
}

func TestTestUtil(testObj *testing.T) {
	suite.Run(testObj, &TestTokenStruct{})
}
func (Ts *TestUtil) TestRandomNumber() {
	num, err := Ts.ut.RandomNumber(0)
	Ts.Require().NotNil(err)
	num, err = Ts.ut.RandomNumber(10000000000000000)
	Ts.Require().Nil(err)
	fmt.Println(num)
}

func (Ts *TestUtil) TestRandomString() {
	st, err := Ts.ut.RandomString(0)
	Ts.Require().NotNil(err)
	st, err = Ts.ut.RandomString(10000000000000000)
	Ts.Require().NotNil(err)
	st, err = Ts.ut.RandomString(10)
	Ts.Require().Nil(err)
	Ts.Require().Equal(len(st) >= 1, true)
	fmt.Println(st)
}

func (Ts *TestUtil) SetupSuite() {
	Ts.ut = Utils{}
}

func (Ts *TestUtil) TearDownSuite() {
	//Ts.ut = nil
}
