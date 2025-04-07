package auth

import "github.com/golang-jwt/jwt/v5"

type TestAuthenticator struct {}


func (t *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	return "", nil
}

func (t *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return nil, nil
}
