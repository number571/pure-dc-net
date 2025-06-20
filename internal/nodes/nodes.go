package nodes

import (
	"crypto/pbkdf2"
	"crypto/sha512"
	"encoding/json"
	"os"

	"github.com/number571/pure-dc-net/pkg/dc"
)

func NodesMapToGenerators(nodes NodesMap) []dc.IGenerator {
	generators := make([]dc.IGenerator, 0, len(nodes))
	for _, n := range nodes {
		generators = append(generators, dc.NewHGenerator(n.GetEncrKey()))
	}
	return generators
}

func LoadNodesMapFromFile(filename string) NodesMap {
	const keySize = 32
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	nodes := make(map[string]*nodeConn, 128)
	if err := json.Unmarshal(body, &nodes); err != nil {
		panic(err)
	}
	resMap := make(NodesMap, len(nodes))
	for k, v := range nodes {
		keyBytes, err := pbkdf2.Key(sha512.New, v.Pasw, []byte("auth_encr"), (1 << 20), 2*keySize)
		if err != nil {
			panic(err)
		}
		v.authKey = keyBytes[:keySize]
		v.encrKey = keyBytes[keySize:]
		resMap[k] = v
	}
	return resMap
}
