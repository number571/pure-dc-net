package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/number571/pure-dc-net/internal/nodes"
)

const (
	dcKeysFile = "dc_keys.json"
)

func loadDCNodesMap() nodes.Nodes {
	filename := filepath.Join(servicePath, dcKeysFile)
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	result := nodes.Nodes{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	return result
}
