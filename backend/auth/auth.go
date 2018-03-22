package auth

import (
	"fmt"
	"log"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/config"
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

func NewTokenManager(method SignMethod) *TokenManager {
	t := TokenManager{}
	switch method {
	case SignMethodSHA256:
		t.algorithm = jwt.SigningMethodHS256
	case SignMethodSHA512:
		t.algorithm = jwt.SigningMethodHS512
	default:
		log.Println("method not valid, using SignMethodSHA512")
		t.algorithm = jwt.SigningMethodHS512
	}
	t.secret = []byte(config.GetHashConfig().Secret)
	return &t
}

func (t *TokenManager) GetToken(inumber string, exp time.Duration) (*models.TokenModel, error) {
	log.Printf("New GetToken request from %s\n", inumber)
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

func (t *TokenManager) GetClaims(token string) (jwt.MapClaims, error) {
	tok, err := t.decodeToken(token)
	if err != nil {
		log.Println("token parse error", err.Error())
		return nil, err
	}

	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("could not parse token")
}

func (t *TokenManager) IsValidToken(token string) error {
	_, err := t.GetClaims(token)
	return err
}

func (t *TokenManager) decodeToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.secret, nil
	})
}
