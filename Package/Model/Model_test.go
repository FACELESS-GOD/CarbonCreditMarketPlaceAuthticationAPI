package Model

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/DevMode"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Utility"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type TestModelStruct struct {
	Ut  Utility.Utils
	Mdl ModelStruct
	suite.Suite
	VerEmail string
	VerPass  string
	UserID   int
}

func (Ts *TestModelStruct) Reset() {
	Ts.VerEmail = ""
	Ts.VerPass = ""
	Ts.UserID = 0

}

func TestModelstructmain(testObj *testing.T) {
	suite.Run(testObj, &TestModelStruct{})
}

func (Ts *TestModelStruct) SetupSuite() {

	util, err := Utility.NewUtility(DevMode.Test, "../../")
	if err != nil {
		Ts.FailNow(err.Error())
	}
	Ts.Ut = util

	mdl, err := NewModel(Ts.Ut)
	if err != nil {
		Ts.FailNow(err.Error())
	}
	Ts.Mdl = mdl
}

func (Its *TestModelStruct) TestVerifyCred() {
	auth := LoginStruct{
		Email:    "",
		Password: "",
	}
	_, _, _, err := Its.Mdl.VerifyCred(auth)
	Its.Require().NotNil(err)

	auth.Email = ""
	auth.Password = Its.VerPass

	_, _, _, err = Its.Mdl.VerifyCred(auth)
	Its.Require().NotNil(err)

	auth.Email = Its.VerEmail
	auth.Password = ""

	_, _, _, err = Its.Mdl.VerifyCred(auth)
	Its.Require().NotNil(err)

	auth.Email = Its.VerEmail
	auth.Password = "q"

	_, _, _, err = Its.Mdl.VerifyCred(auth)
	Its.Require().NotNil(err)

	auth.Email = "q"
	auth.Password = Its.VerPass

	_, _, _, err = Its.Mdl.VerifyCred(auth)
	Its.Require().NotNil(err)

	auth.Email = Its.VerEmail
	auth.Password = Its.VerPass

	IsValid, refereshToken, accessToken, err := Its.Mdl.VerifyCred(auth)
	Its.Require().Nil(err)
	Its.Require().Equal(IsValid, true)
	Its.Require().Equal(len(refereshToken) >= 1, true)
	Its.Require().Equal(len(accessToken) >= 1, true)

}

func (Its *TestModelStruct) TestAddRefereshTokenToDB() {

	_, _, err := Its.Mdl.AddRefereshTokenToDB(0)
	Its.Require().NotNil(err)

	refereshToken, accessToken, err := Its.Mdl.AddRefereshTokenToDB(Its.UserID)

	Its.Require().Nil(err)
	Its.Require().Equal(len(refereshToken) >= 1, true)
	Its.Require().Equal(len(accessToken) >= 1, true)

}

func (Its *TestModelStruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestVerifyCred":
		Its.Reset()
		randNum, err := Its.Ut.RandomNumber(110)
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
	case "TestAddRefereshTokenToDB":
		Its.Reset()
		randNum, err := Its.Ut.RandomNumber(110)
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

	}
}

func (Its *TestModelStruct) AfterTest(SuiteName string, TestName string) {

	switch TestName {
	case "TestVerifyCred":
		err := Its.deleteUser()

		if err != nil {
			Its.FailNow(err.Error())
		}
	case "TestAddRefereshTokenToDB":

		err := Its.deleteUser()

		if err != nil {
			Its.FailNow(err.Error())
		}
	}
}

func (Its *TestModelStruct) TearDownSuite() {
	Its.Ut.DB.Close()
}

func (Its *TestModelStruct) GenerateHash(Password string) (string, error) {

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

func (Its *TestModelStruct) addUser() error {

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

func (Its *TestModelStruct) deleteUser() error {

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
