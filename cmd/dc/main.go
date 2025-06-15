package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/number571/pure-dc-net/internal/nodes"
	"github.com/number571/pure-dc-net/internal/service"
	"github.com/number571/pure-dc-net/internal/token"
	"github.com/number571/pure-dc-net/pkg/dc"
)

var (
	servicePath = os.Getenv("SERVICE_PATH")
	consumeAddr = os.Getenv("CONSUME_ADDR")
	produceAddr = os.Getenv("PRODUCE_ADDR")
)

var (
	currIter = loadDCIter()
	nodeName = loadDCName()
	nodesMap = loadDCNodesMap()
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var (
		dcNet  = dc.NewDCNet(currIter, nodes.NodesKeysToGenerators(nodesMap)...)
		ttlzr  = dc.NewTotalizer()
		bqueue = make(chan byte, 512)
	)

	log.Println("service is listening...")

	go runGenerator(ctx, dcNet, ttlzr, bqueue)

	go func() {
		server := service.NewDCProducerServer(produceAddr, bqueue)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	server := service.NewDCConsumerServer(consumeAddr, nodesMap, dcNet, ttlzr)
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
			b := byte(0)
			select {
			case x := <-bqueue:
				b = x
			default:
			}
			gb := b ^ dcNet.Generate()
			doDCRequest(ctx, dcNet, nodesMap, nodeName, gb)
			for {
				if ttlzr.Size() == len(nodesMap) {
					break
				}
				time.Sleep(time.Second)
			}
			ttlzr.Store(gb)
			if r := ttlzr.Sum(); r != 0 {
				log.Println(string(r))
			}
		}
	}
}

func doDCRequest(ctx context.Context, dcNet dc.IDCNet, nodes nodes.Nodes, nname string, gb byte) {
	defer func() { storeDCIter(dcNet.Iteration()) }()

	tokenData := token.MarshalTokenData(&token.TokenData{
		Name: nname,
		Iter: dcNet.Iteration(),
		Byte: gb,
	})

	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))

	for _, node := range nodes {
		node := node
		go func() {
			defer wg.Done()

			token := token.GenerateToken(node.Key, tokenData)
			_ = service.ConsumeRequest(ctx, node.Addr, token)
		}()
	}

	wg.Wait()
}
