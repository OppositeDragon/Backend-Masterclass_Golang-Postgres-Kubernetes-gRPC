package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrKeySizeTooSmall
	}

	return &JwtMaker{secretKey}, nil
}

func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// / VerifyToken checks if the token is valid or not
func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, isOk := token.Method.(*jwt.SigningMethodHMAC)
		if !isOk {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		vErr, isOk := err.(*jwt.ValidationError)
		if isOk && errors.Is(vErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken

	}
	payload, isOk := jwtToken.Claims.(*Payload)
	if !isOk {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
