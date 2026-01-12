package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJwtAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (authenticator *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(authenticator.secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authenticator *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(authenticator.secret), nil
	}, jwt.WithExpirationRequired(),
		jwt.WithAudience(authenticator.aud),
		jwt.WithIssuer(authenticator.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	return token, err
}
