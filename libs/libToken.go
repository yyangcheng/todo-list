package libs

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/yyangc/todo-list/config"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

func CreateToken(userId uint64) (*TokenDetails, error) {
	td := new(TokenDetails)
	td.AtExpires = time.Now().Add(time.Hour * 1).Unix()
	accessUuid, _ := uuid.NewV4()
	td.AccessUuid = accessUuid.String()

	var err error
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"accessUuid": td.AccessUuid,
		"userId":     userId,
		"exp":        td.AtExpires,
	})
	// Sign and get the complete encoded token as a string using the secret
	td.AccessToken, err = at.SignedString([]byte(config.Env.JWT.AccessSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["accessUuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)

		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.JWT.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
