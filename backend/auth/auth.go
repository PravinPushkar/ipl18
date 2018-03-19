package auth

import (
	"log"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/models"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	algorithm jwt.SigningMethod
	secret    []byte
}

const (
	SignMethodSHA256 = iota << 1
	SignMethodSHA512
)

type SignMethod int
type ValidationMethod func(string) bool

func NewTokenManager(method SignMethod, secret string) *TokenManager {
	t := TokenManager{}
	switch method {
	case SignMethodSHA256:
		t.algorithm = jwt.SigningMethodRS256
	case SignMethodSHA512:
		t.algorithm = jwt.SigningMethodHS512
	default:
		log.Println("method not valid, using SignMethodSHA512")
		t.algorithm = jwt.SigningMethodHS512
	}
	t.secret = []byte(secret)
	return &t
}

func (t *TokenManager) GetToken(inumber string, exp time.Duration) (*models.TokenModel, error) {
	token := jwt.New(t.algorithm)
	claims := token.Claims.(jwt.MapClaims)

	claims["inumber"] = inumber
	claims["exp"] = time.Now().Add(time.Hour * exp).Unix()

	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	tModel := models.TokenModel{}
	tModel.Token = tokenString

	return &tModel, nil
}

func (t *TokenManager) IsValidToken(method ValidationMethod) bool {
	return true
}

func (t *TokenManager) ParseToken(token string) bool {
	return true
}
