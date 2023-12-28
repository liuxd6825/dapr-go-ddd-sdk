package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

var defaultJWTKey = []byte("#@!{[duXm-serVice-t0ken]},.(10086)$!")

type loginClaim struct {
	User *user
	jwt.RegisteredClaims
}

func decodeJwt(jwtText string, jwtKey string) (LoginUser, error) {
	if len(jwtText) <= 6 {
		return nil, errors.New("jwtText长度不足6")
	}
	var newJwtString string

	tokenType := jwtText[0:6]
	if strings.ToLower(tokenType) == "bearer" {
		newJwtString = jwtText[7:len(jwtText)]
	} else {
		newJwtString = jwtText
	}

	token, err := jwt.ParseWithClaims(newJwtString, &loginClaim{}, func(token *jwt.Token) (interface{}, error) {
		if jwtKey == "" {
			return defaultJWTKey, nil
		}
		return []byte(jwtKey), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*loginClaim)
	if !ok {
		return nil, errors.New("非法的Token内容")
	}
	return claims.User, nil
}
