package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySecret = []byte("howoriginal")
var tokenExpireTime = time.Hour * 24

type Claims struct {
	UserID int64 `json:"user_id"`
	Username string `json:"username"`

	jwt.StandardClaims
}


func GenToken(userID int64, username string) (string, error) {
	c := Claims{
		UserID: userID,
		Username: username,
		
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireTime).Unix(),
			Issuer: "middit",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(mySecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	mc := Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, &mc, func(t *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return &mc, nil
	}

	return nil, errors.New("invalid token")
}