package main

import (
	"context"
	"log"
	"os"
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

var (
	nodesMap = loadDCNodesMap()
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var (
		dcNet  = dc.NewDCNet(0, nodes.NodesKeysToGenerators(nodesMap)...)
		ttlzr  = dc.NewTotalizer()
		bqueue = make(chan byte, 512)
	)

	log.Println("service is listening...")

	go runGenerator(ctx, dcNet, ttlzr, bqueue)

	go func() {
		server := service.NewDCInternalServer(produceAddr, bqueue)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	server := service.NewDCExternalServer(consumeAddr, nodesMap, dcNet, ttlzr)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func runGenerator(ctx context.Context, dcNet dc.IDCNet, ttlzr dc.ITotalizer, bqueue chan byte) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b := dcNet.Generate()
			select {
			case x := <-bqueue:
				b ^= x // send message
			default:
			}
			iteration := dcNet.Iteration()
			tokenData := token.MarshalTokenData(&token.TokenData{
				Name: serviceName,
				Iter: iteration,
				Byte: b,
			})
			service.DoRequest(ctx, nodesMap, tokenData)
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
