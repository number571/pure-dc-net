package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/pure-dc-net/internal/nodes"
)

const (
	dcIterFile = "dc_iter.txt"
	dcKeysFile = "dc_keys.txt"
	dcNameFile = "dc_name.txt"
)

func storeDCIter(iter uint64) {
	filename := filepath.Join(servicePath, dcIterFile)
	err := os.WriteFile(filename, []byte(strconv.FormatUint(iter, 10)), 0600)
	if err != nil {
		panic(err)
	}
}

func loadDCIter() uint64 {
	filename := filepath.Join(servicePath, dcIterFile)
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

func loadDCName() string {
	filename := filepath.Join(servicePath, dcNameFile)
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func loadDCNodesMap() nodes.Nodes {
	filename := filepath.Join(servicePath, dcKeysFile)
	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	mapping := bytes.Split(body, []byte("\n"))
	result := make(nodes.Nodes, len(mapping))
	for _, m := range mapping {
		r := bytes.Split(m, []byte(";"))
		if len(r) < 3 {
			panic("invalid mapping")
		}
		result[string(r[0])] = &nodes.NodeConn{
			Addr: string(r[1]),
			SKey: bytes.Join(r[2:], []byte{}),
		}
	}
	return result
}
