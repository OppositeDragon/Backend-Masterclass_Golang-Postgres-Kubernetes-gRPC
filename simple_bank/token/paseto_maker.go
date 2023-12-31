package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto *paseto.V2
	key    []byte
}

func NewPasetoMaker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be exatly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		key:    []byte(key),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(username, duration)
	token, err := maker.paseto.Encrypt(maker.key, payload, []byte("optionalfooter"))
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.key, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
