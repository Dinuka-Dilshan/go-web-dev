package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJwtAuthenticator(secret, aud, iss string) JWTAuthenticator {
	return JWTAuthenticator{
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

func (*JWTAuthenticator) ValidateToken() (*jwt.Token, error) {
	return nil, nil
}
