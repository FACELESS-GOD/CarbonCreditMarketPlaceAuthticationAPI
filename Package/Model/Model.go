package Model

import (
	"context"
	"errors"
	"time"

	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Helper/TokenType"
	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Utility"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginStruct struct {
	Email    string
	Password string
}

type UserToken struct {
	UserID   int
	TokenID  uuid.UUID
	RoleId   int
	IssuedAT time.Time
	jwt.RegisteredClaims
}

type ModelStruct struct {
	Ut            Utility.Utils
	IsAnyError    bool
	ErrorMessages []string
}

func (Mdl *ModelStruct) Reset() {
	Mdl.IsAnyError = false
	Mdl.ErrorMessages = make([]string, 0)
}

const GetUserIDQuery string = `
SELECT UserId FROM User
WHERE Is_Visible = true AND email = ?
LIMIT 1
;
`
const GetHashQuery string = `
SELECT Hash_Password FROM UserCred
WHERE Is_Visible = true AND UserId = ?
LIMIT 1
;
`

func NewModel(Ut Utility.Utils) (ModelStruct, error) {
	mdl := ModelStruct{}
	mdl.Ut = Ut
	return mdl, nil
}

func (Mdl *ModelStruct) VerifyCred(AuthDt LoginStruct) (bool, string, string, error) {
	Mdl.Reset()

	if len(AuthDt.Email) < 1 {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Invalid Email")
		return false, "", "", errors.New("Invalid Email")
	}

	if len(AuthDt.Password) < 1 {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Invalid Length of Password")
		return false, "", "", errors.New("Invalid Length of Password")
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Ut.DB.BeginTx(ctx, &Mdl.Ut.TxOption)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, "", "", err
	}

	response, err := db.Query(GetUserIDQuery, AuthDt.Email)

	if err != nil {
		response.Close()
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+","+nerr.Error())
			return false, "", "", errors.New(err.Error() + "," + nerr.Error())
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, "", "", err
		}
	}

	var userID int = 0

	for response.Next() {
		err := response.Scan(&userID)

		if err != nil {
			response.Close()
			nerr := db.Rollback()
			if nerr != nil {
				Mdl.IsAnyError = true
				Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+","+nerr.Error())
				return false, "", "", errors.New(err.Error() + "," + nerr.Error())
			} else {
				Mdl.IsAnyError = true
				Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
				return false, "", "", err
			}
		}
	}

	if userID < 1 {
		response.Close()
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Email Does Not exits ,"+nerr.Error())
			return false, "", "", errors.New("Email Does Not exits" + "," + nerr.Error())
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Email Does Not exits")
			return false, "", "", errors.New("Email Does Not exits")
		}
	}

	response.Close()

	response, err = db.Query(GetHashQuery, userID)

	if err != nil {
		response.Close()
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+","+nerr.Error())
			return false, "", "", errors.New(err.Error() + "," + nerr.Error())
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, "", "", err
		}
	}

	var hashPassword string = ""

	for response.Next() {
		err := response.Scan(&hashPassword)

		if err != nil {
			response.Close()
			nerr := db.Rollback()
			if nerr != nil {
				Mdl.IsAnyError = true
				Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+","+nerr.Error())
				return false, "", "", errors.New(err.Error() + "," + nerr.Error())
			} else {
				Mdl.IsAnyError = true
				Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
				return false, "", "", err
			}
		}

	}

	if len(hashPassword) < 1 {
		response.Close()
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Password is Not Present ,"+nerr.Error())
			return false, "", "", errors.New("Password is Not Present" + "," + nerr.Error())
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Password is Not Present")
			return false, "", "", errors.New("Password is Not Present")
		}
	}

	response.Close()

	err = db.Rollback()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, "", "", errors.New(err.Error())
	}

	isValid, err := Mdl.verifyHashPassword(hashPassword, AuthDt.Password)

	if isValid != true {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Password is Not Correct")
		return false, "", "", errors.New("Password is Not Correct")
	}

	tkn, err := Mdl.Ut.CreateToken(userID, TokenType.RefereshToken)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, "", "", errors.New(err.Error())
	}

	accessTkn, err := Mdl.Ut.CreateToken(userID, TokenType.AccessToken)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, "", "", errors.New(err.Error())
	}

	return true, tkn, accessTkn, nil

}

const deleteToken string = `
UPDATE TokenStore
SET Is_Visible = 0 , Last_Modified_Date = CURDATE()
WHERE UserId = ?
;
`

const InsertToken string = `
INSERT INTO TokenStore (
  UserId, Token
) VALUES (
  ? , ? 
)
;
`

func (Mdl *ModelStruct) AddRefereshTokenToDB(UserID int) (string, string, error) {
	if UserID < 1 {
		return "", "", errors.New("Invalid Data")
	}

	newtkn, err := Mdl.Ut.CreateToken(UserID, TokenType.RefereshToken)
	if err != nil {
		return "", "", err
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Ut.DB.BeginTx(ctx, &Mdl.Ut.TxOption)
	if err != nil {
		return "", "", err
	}

	response, err := db.Query(deleteToken, UserID)

	if err != nil {
		response.Close()
		return "", "", err
	}
	response.Close()

	response, err = db.Query(InsertToken, UserID, newtkn)

	if err != nil {
		response.Close()
		return "", "", err
	}
	response.Close()

	err = db.Rollback()
	if err != nil {
		return "", "", err
	}

	accesstkn, err := Mdl.Ut.CreateToken(UserID, TokenType.AccessToken)
	if err != nil {
		return "", "", err
	}

	return newtkn, accesstkn, err

}

func (Ut *ModelStruct) verifyHashPassword(Hash string, Pass string) (bool, error) {
	if len(Hash) < 1 || len(Pass) < 1 {
		return false, errors.New("Invalid Data")
	}

	err := bcrypt.CompareHashAndPassword([]byte(Hash), []byte(Pass))

	if err != nil {
		return false, err
	}

	return true, nil

}
