package utils

import (
	"encoding/base64"
	"errors"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type jwtClaims struct {
	jwt.StandardClaims
	Id        string
	Username  string
	Role      string
	TokenType string
	ApiCode   string
}

func (w *JwtWrapper) GenerateToken(user *domain.User) (signedToken string, err error) {
	claims := &jwtClaims{
		Id:        user.Id.Hex(),
		Username:  user.Username,
		Role:      domain.ConvertRoleToString(user.Role),
		TokenType: "JWT",
		ApiCode:   "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(8)).Unix(),
			Issuer:    w.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, _ := base64.URLEncoding.DecodeString("dislinkt")
	signedToken, err = token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *JwtWrapper) ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return base64.URLEncoding.DecodeString("dislinkt")
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil

}
