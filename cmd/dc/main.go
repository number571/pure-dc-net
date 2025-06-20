package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/number571/pure-dc-net/internal/nodes"
	"github.com/number571/pure-dc-net/internal/service"
	"github.com/number571/pure-dc-net/internal/token"
	"github.com/number571/pure-dc-net/pkg/dc"
)

var (
	serviceName = os.Getenv("SERVICE_NAME")
	servicePath = os.Getenv("SERVICE_PATH")
	consumeAddr = os.Getenv("CONSUME_ADDR")
	produceAddr = os.Getenv("PRODUCE_ADDR")
)

func main() {
	nodesMap := nodes.LoadNodesMapFromFile(filepath.Join(servicePath, "dc.json"))

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var (
		dcState = dc.NewDCState(0, nodes.NodesMapToGenerators(nodesMap)...)
		ttlzr   = dc.NewTotalizer()
		bqueue  = make(chan byte, 512)
	)

	log.Println("service is listening...")

	go runGenerator(ctx, nodesMap, dcState, ttlzr, bqueue)

	go func() {
		server := service.NewDCInternalServer(produceAddr, bqueue)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	server := service.NewDCExternalServer(consumeAddr, nodesMap, dcState, ttlzr)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func runGenerator(ctx context.Context, nodesMap nodes.NodesMap, dcState dc.IDCState, ttlzr dc.ITotalizer, bqueue chan byte) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			i, b := generate(dcState, bqueue)
			tokenData := &token.TokenData{
				Name: serviceName,
				Iter: i,
				Byte: b,
			}
			service.Commit(ctx, nodesMap, tokenData)
			for {
				if ttlzr.Size() == len(nodesMap) {
					break
				}
				time.Sleep(time.Second)
			}
			ttlzr.Store(b)
			if r := ttlzr.Sum(); r != 0 {
				log.Println(string(r))
			}
		}
	}
}

func generate(dcState dc.IDCState, bqueue chan byte) (uint64, byte) {
	b := dcState.Generate()
	select {
	case x := <-bqueue:
		b ^= x // push message
	default:
	}
	return dcState.Iteration(), b
}
