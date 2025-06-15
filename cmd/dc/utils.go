package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/number571/pure-dc-net/pkg/dc"
)

func sumGenerated(gb map[string]byte) byte {
	s := byte(0)
	for _, b := range gb {
		s ^= b
	}
	return s
}

func generateToken(keys map[string][]byte, myAddr, dstAddr string) *reqToken {
	salt := make([]byte, 32)
	return &reqToken{
		Addr: myAddr,
		Salt: salt,
		Mac:  hmac.New(sha512.New, keys[dstAddr]).Sum(salt),
	}
}

func validateMAC(keys map[string][]byte, token *reqToken) error {
	key, ok := keys[token.Addr]
	if !ok {
		return errors.New("unknown address")
	}
	vmac := hmac.New(sha512.New, key).Sum(token.Salt)
	if !bytes.Equal(token.Mac, vmac) {
		return errors.New("invalid mac")
	}
	return nil
}

func keysToGenerators(keys map[string][]byte) []dc.IGenerator {
	generators := make([]dc.IGenerator, 0, len(keys))
	for _, k := range keys {
		generators = append(generators, dc.NewHGenerator(k))
	}
	return generators
}

func storeDCIter(filename string, iter uint64) {
	err := os.WriteFile(filename, []byte(strconv.FormatUint(iter, 10)), 0600)
	if err != nil {
		panic(err)
	}
}

func loadDCAddr(filename string) string {
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func loadDCIter(filename string) uint64 {
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	iter, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		panic(err)
	}
	return iter
}

func loadDCKeys(filename string) map[string][]byte {
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	mapping := bytes.Split(body, []byte("\n"))
	result := make(map[string][]byte, len(mapping))
	for _, m := range mapping {
		r := bytes.Split(m, []byte(";"))
		if len(r) != 2 {
			fmt.Println(string(m))
			panic("invalid mapping")
		}
		result[string(r[0])] = r[1]
	}
	return result
}
