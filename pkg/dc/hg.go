package dc

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"hash"
)

type hGenerator struct {
	h hash.Hash
}

func NewHGenerator(k []byte) IGenerator {
	return &hGenerator{
		h: hmac.New(sha512.New, k),
	}
}

func (p *hGenerator) Generate(i uint64) byte {
	iter := make([]byte, 8)
	binary.BigEndian.PutUint64(iter, i)
	return p.h.Sum(iter)[0]
}
