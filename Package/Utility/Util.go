package Utility

import (
	"database/sql"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/FACELESS-GOD/CarbonCreditMarketPlaceAuthticationAPI/Package/Configurator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ===========  LOCAL ENTITIES ==========

func (Ut *Utils) initiate() error {
	err := Ut.initiateDB()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (Ut *Utils) initiateDB() error {

	db, err := sql.Open(Ut.config.DBDRIVER, Ut.config.DBCONNSTRING)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	Ut.DB = db

	txOption := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	Ut.TxOption = txOption

	return nil
}

// ==========================================================================================================================================

// ========== PUBLIC ENTITIES ==========

type Utils struct {
	config    Configurator.Config
	DB        *sql.DB
	RDB       *redis.Client
	rdbOption redis.Options
	Mode      int
	TxOption  sql.TxOptions
}

type TokenMetaData struct {
	TokenID   uuid.UUID `json:"tokenid"`
	UserID    int       `json:"userid"`
	IssuedAT  time.Time `json:"issuedat"`
	ExpiredAt time.Time `json:"expiredat"`
	RoleID    int       `json:"roleid"`
	jwt.RegisteredClaims
}

func NewUtility(Mode int, Path string) (Utils, error) {
	ut := Utils{}
	ut.Mode = Mode
	conf, err := Configurator.NewConfigurator(Path)
	if err != nil {
		log.Println(err)
		return ut, err
	}

	ut.config = conf

	err = ut.initiate()
	if err != nil {
		log.Println(err)
		return ut, err
	}

	return ut, nil
}

func (Ut *Utils) CreateToken(userId int, TokenType string) (string, error) {

	if userId < 1 {
		return "", errors.New("Invalid Data!.")
	}

	randomID, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	tkm := TokenMetaData{
		TokenID:   randomID,
		UserID:    userId,
		IssuedAT:  time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(time.Now().Day())),
	}

	if TokenType != "Referesh" {
		tkm.ExpiredAt = time.Now().Add(time.Hour)
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"TokenID":   tkm.TokenID,
		"UserId":    tkm.UserID,
		"IssuedAT":  tkm.IssuedAT,
		"ExpiredAT": tkm.ExpiredAt,
	})

	tokenString, err := tkn.SignedString([]byte(Ut.config.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (Ut *Utils) VerifyToken(Token string) (bool, TokenMetaData, error) {
	tkm := TokenMetaData{}

	tkn, err := jwt.ParseWithClaims(Token, &tkm, func(token *jwt.Token) (interface{}, error) {
		return []byte(Ut.config.JwtSecretKey), nil
	}, jwt.WithLeeway(5*time.Hour))

	if err != nil {
		return false, tkm, err
	}

	if tkn.Valid != true {
		return false, tkm, err
	}
	return true, tkm, err
}

// Generate Randomnumber range between 0 and Salt .
func (Ut *Utils) RandomNumber(Salt int) (int32, error) {
	if Salt < 1 {
		return 0, errors.New("RandomSalt is less than 1")
	}
	return rand.Int32N(int32(Salt)), nil
}

// Generate Random String of lenght num and num cannot be less than 0 and greater than 256.
func (Ut *Utils) RandomString(num int) (string, error) {

	if num < 1 {
		return "", errors.New("Length of string cannot be less than 1")
	}

	if num > 254 {
		return "", errors.New("Length of string cannot be greated than 256")
	}

	base := "QWERTYUIOPASDFGHJKLZXCVBNM1234567890"

	baseRandomString := []byte{}

	for i := 1; i <= num; i++ {
		newNumIndex, err := Ut.RandomNumber(len(base))
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		baseRandomString = append(baseRandomString, base[newNumIndex])
	}

	return string(baseRandomString), nil

}

// =========================================================================================================================================
