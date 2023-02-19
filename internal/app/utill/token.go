package utill

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

type TokenJWT struct {
	signingKey     []byte
	ExpiryInterval time.Duration
}

func NewToken(signingKey string, expiryInterval time.Duration) *TokenJWT {
	return &TokenJWT{signingKey: []byte(signingKey), ExpiryInterval: expiryInterval}
}

func (t *TokenJWT) GenerateToken(sub int) (string, error) {
	return t.generateToken(strconv.Itoa(sub))
}

func (t *TokenJWT) ParseToken(tokenStr string) (int, error) {
	subStr, err := t.parseToken(tokenStr)
	if err != nil {
		return 0, err
	}
	sub, err := strconv.Atoi(subStr)
	if err != nil {
		return 0, err
	}
	return sub, nil
}

func (t *TokenJWT) generateToken(sub string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(t.ExpiryInterval).Unix(),
		Subject:   sub,
	})
	return token.SignedString(t.signingKey)
}

func (t *TokenJWT) parseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return t.signingKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", errors.New("token claims are not type of  jwt.StandardClaims")
	}
	return claims.Subject, nil
}
