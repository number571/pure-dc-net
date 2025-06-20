package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
)

type Token struct {
	Data []byte `json:"data"`
	Salt []byte `json:"salt"`
	Mac  []byte `json:"mac"`
}

func GenerateToken(key, data []byte) *Token {
	salt := make([]byte, 32)
	return &Token{
		Salt: salt,
		Data: data,
		Mac: hmac.New(sha512.New, key).Sum(
			bytes.Join([][]byte{salt, data}, []byte{}),
		),
	}
}

func ValidateMAC(key []byte, token *Token) error {
	vmac := hmac.New(sha512.New, key).Sum(
		bytes.Join([][]byte{token.Salt, token.Data}, []byte{}),
	)
	if !bytes.Equal(token.Mac, vmac) {
		return errors.New("invalid mac")
	}
	return nil
}
