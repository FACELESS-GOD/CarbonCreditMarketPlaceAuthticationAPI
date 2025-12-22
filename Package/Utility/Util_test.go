package Utility

import (
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
