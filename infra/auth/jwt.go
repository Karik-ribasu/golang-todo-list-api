package auth

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthToken struct {
	Token     string `json:"token,omitempty"`
	IssuedAt  int64  `json:"issued_at,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

type CustomClaims struct {
	UserUUID string `json:"user_uuid"`
	jwt.StandardClaims
}

const expiration = time.Hour * 24
const issuer = "todo-list-api"

func GenerateJWTToken(privateKey *rsa.PrivateKey, userUUID string) (token AuthToken, err error) {

	claims := CustomClaims{}
	claims.UserUUID = userUUID
	claims.ExpiresAt = time.Now().Add(expiration).Unix()
	claims.IssuedAt = time.Now().Unix()
	claims.Issuer = issuer

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := jwtToken.SignedString(privateKey)

	authToken := AuthToken{
		Token:     tokenString,
		ExpiresAt: claims.ExpiresAt,
		IssuedAt:  claims.IssuedAt,
	}

	return authToken, err
}

func ParseAndValidateJWTtoken(privateKey *rsa.PrivateKey, tokenString string) (claims CustomClaims, isValid bool, err error) {

	claims = CustomClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(jwt *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return claims, false, err
	}

	return claims, true, nil
}
