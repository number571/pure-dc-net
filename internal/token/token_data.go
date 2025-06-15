package token

import (
	"encoding/json"
	"fmt"
)

type TokenData struct {
	Name string `json:"name"`
	Iter uint64 `json:"iter"`
	Byte byte   `json:"byte"`
}

func MarshalTokenData(tokenData *TokenData) []byte {
	res, err := json.Marshal(tokenData)
	if err != nil {
		panic(err)
	}
	return res
}

func UnmarshalTokenData(bytesData []byte) (*TokenData, error) {
	tokenData := &TokenData{}
	if err := json.Unmarshal(bytesData, tokenData); err != nil {
		fmt.Println(string(bytesData))
		return nil, err
	}
	return tokenData, nil
}
